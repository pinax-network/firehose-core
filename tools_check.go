// Copyright 2021 dfuse Platform Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package firecore

import (
	"fmt"
	"os"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
	"github.com/streamingfast/bstream"
	"github.com/streamingfast/cli"
	"github.com/streamingfast/cli/sflags"
	"github.com/streamingfast/dstore"
	"github.com/streamingfast/firehose-core/firehose/tools"
	pbbstream "github.com/streamingfast/pbgo/sf/bstream/v1"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

var toolsCheckCmd = &cobra.Command{Use: "check", Short: "Various checks for deployment, data integrity & debugging"}

var toolsCheckForksCmd = &cobra.Command{
	Use:   "forks <forked-blocks-store-url>",
	Short: "Reads all forked blocks you have and print longest linkable segments for each fork",
	Args:  cobra.ExactArgs(1),
}

var toolsCheckMergedBlocksCmd = &cobra.Command{
	// TODO: Not sure, it's now a required thing, but we could probably use the same logic as `start`
	//       and avoid altogether passing the args. If this would also load the config and everything else,
	//       that would be much more seamless!
	Use:   "merged-blocks <store-url>",
	Short: "Checks for any holes in merged blocks as well as ensuring merged blocks integrity",
	Args:  cobra.ExactArgs(1),
}

func init() {
	toolsCmd.AddCommand(toolsCheckCmd)
	toolsCheckCmd.AddCommand(toolsCheckForksCmd)
	toolsCheckCmd.AddCommand(toolsCheckMergedBlocksCmd)

	toolsCheckCmd.PersistentFlags().StringP("range", "r", "", "Block range to use for the check")

	toolsCheckMergedBlocksCmd.Flags().BoolP("print-stats", "s", false, "Natively decode each block in the segment and print statistics about it, ensuring it contains the required blocks")
	toolsCheckMergedBlocksCmd.Flags().BoolP("print-full", "f", false, "Natively decode each block and print the full JSON representation of the block, should be used with a small range only if you don't want to be overwhelmed")

	toolsCheckForksCmd.Flags().Uint64("min-depth", 1, "Only show forks that are at least this deep")
	toolsCheckForksCmd.Flags().Uint64("after-block", 0, "Only show forks that happened after this block number, if value is not 0")
}

func configureToolsCheckCmd[B Block](chain *Chain[B]) {
	blockPrinter := chain.BlockPrinter()

	toolsCheckMergedBlocksCmd.RunE = createToolsCheckMergedBlocksE(blockPrinter)
	toolsCheckMergedBlocksCmd.Example = ExamplePrefixed(chain, "tools check merged-blocks", `
		"./sf-data/storage/merged-blocks"
		"gs://<project>/<bucket>/<path>" -s
		"s3://<project>/<bucket>/<path>" -f
		"az://<project>/<bucket>/<path>" -r "10 000 - 1 000 000"
	`)

	toolsCheckForksCmd.RunE = toolsCheckForksE
}

func createToolsCheckMergedBlocksE(blockPrinter BlockPrinterFunc) CommandExecutor {
	return func(cmd *cobra.Command, args []string) error {
		storeURL := args[0]
		fileBlockSize := uint64(100)

		blockRange, err := tools.GetBlockRangeFromFlag(cmd, "range")
		if err != nil {
			return err
		}

		printDetails := tools.PrintNoDetails
		if sflags.MustGetBool(cmd, "print-stats") {
			printDetails = tools.PrintStats
		}

		if sflags.MustGetBool(cmd, "print-full") {
			printDetails = tools.PrintFull
		}

		return tools.CheckMergedBlocks(cmd.Context(), rootLog, storeURL, fileBlockSize, blockRange, func(block *pbbstream.Block) {
			blockPrinter(block, false, os.Stdout)
		}, printDetails)
	}
}

func toolsCheckForksE(cmd *cobra.Command, args []string) error {
	storeURL := args[0]

	blocksStore, err := dstore.NewDBinStore(storeURL)
	cli.NoError(err, "unable to create blocks store")

	oneBlockFiles := []*bstream.OneBlockFile{}
	oneBlockFilesByID := map[string]*bstream.OneBlockFile{}
	err = blocksStore.Walk(cmd.Context(), "", func(filename string) error {
		file, err := bstream.NewOneBlockFile(filename)
		cli.NoError(err, "unable to parse block filename %q", filename)

		oneBlockFiles = append(oneBlockFiles, file)
		oneBlockFilesByID[file.ID] = file
		return nil
	})
	cli.NoError(err, "unable to walk blocks store")

	if len(oneBlockFiles) == 0 {
		fmt.Println("No forked blocks found")
	}

	parentOf := map[string]*bstream.OneBlockFile{}
	for _, file := range oneBlockFiles {
		parentOf[file.ID] = oneBlockFilesByID[file.PreviousID]
	}

	// All forks indexed by their lowest height parent, list is always sorted by height
	links := map[string][]*bstream.OneBlockFile{}
	for _, file := range oneBlockFiles {
		links[file.ID] = []*bstream.OneBlockFile{file}
	}

	for {
		changed := false
		for root, chain := range links {
			if parent := parentOf[root]; parent != nil {
				delete(links, root)
				links[parent.ID] = append([]*bstream.OneBlockFile{parent}, chain...)
				changed = true
			}

		}
		if !changed {
			break
		}
	}

	sortedKeys := maps.Keys(links)
	slices.SortFunc(sortedKeys, func(a, b string) bool {
		return links[a][0].Num < links[b][0].Num
	})

	minDepth := sflags.MustGetInt(cmd, "min-depth")
	afterBlock := sflags.MustGetUint64(cmd, "after-block")

	for _, key := range sortedKeys {
		link := links[key]

		if len(link) < int(minDepth) {
			continue
		}

		if afterBlock != 0 && link[0].Num <= afterBlock {
			continue
		}

		chain := make([]string, len(link))
		for i, segment := range link {
			spaces := strings.Repeat(" ", (i)+2)

			canonical := ""
			if i == 0 {
				canonical = " (on chain)"
			}

			chain[i] = fmt.Sprintf(spaces+"#%d [%s <= %s%s]", blockNumber(segment.Num), segment.ID, segment.PreviousID, canonical)
		}

		fmt.Printf("Fork Depth %d\n%s\n\n", len(link), strings.Join(chain, "\n"))
	}

	return nil
}

type blockNumber uint64

func (b blockNumber) String() string {
	return fmt.Sprintf("#%s", humanize.Comma(int64(b)))
}
