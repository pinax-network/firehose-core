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

package print

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/streamingfast/bstream"
	pbbstream "github.com/streamingfast/bstream/pb/sf/bstream/v1"
	"github.com/streamingfast/cli"
	"github.com/streamingfast/dstore"
	firecore "github.com/streamingfast/firehose-core"
	"github.com/streamingfast/firehose-core/types"
)

func NewToolsPrintCmd[B firecore.Block](chain *firecore.Chain[B]) *cobra.Command {
	toolsPrintCmd := &cobra.Command{
		Use:   "print",
		Short: "Prints of one block or merged blocks file",
	}

	toolsPrintOneBlockCmd := &cobra.Command{
		Use:   "one-block <store> <block_num>",
		Short: "Prints a block from a one-block file",
		Args:  cobra.ExactArgs(2),
	}

	toolsPrintMergedBlocksCmd := &cobra.Command{
		Use:   "merged-blocks <store> <start_block>",
		Short: "Prints the content summary of a merged blocks file.",
		Args:  cobra.ExactArgs(2),
	}

	toolsPrintCmd.AddCommand(toolsPrintOneBlockCmd)
	toolsPrintCmd.AddCommand(toolsPrintMergedBlocksCmd)

	toolsPrintCmd.PersistentFlags().Bool("transactions", false, "When in 'text' output mode, also print transactions summary")

	toolsPrintOneBlockCmd.RunE = createToolsPrintOneBlockE(chain)
	toolsPrintMergedBlocksCmd.RunE = createToolsPrintMergedBlocksE(chain)

	return toolsPrintCmd
}

func createToolsPrintMergedBlocksE[B firecore.Block](chain *firecore.Chain[B]) firecore.CommandExecutor {
	return func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		outputPrinter, err := GetOutputPrinter(cmd, chain.BlockFileDescriptor())
		cli.NoError(err, "Unable to get output printer")

		storeURL := args[0]
		store, err := dstore.NewDBinStore(storeURL)
		if err != nil {
			return fmt.Errorf("unable to create store at path %q: %w", store, err)
		}

		startBlock, err := strconv.ParseUint(args[1], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid base block %q: %w", args[1], err)
		}
		blockBoundary := types.RoundToBundleStartBlock(startBlock, 100)

		filename := fmt.Sprintf("%010d", blockBoundary)
		reader, err := store.OpenObject(ctx, filename)
		if err != nil {
			fmt.Printf("❌ Unable to read blocks filename %s: %s\n", filename, err)
			return err
		}
		defer reader.Close()

		readerFactory, err := bstream.NewDBinBlockReader(reader)
		if err != nil {
			fmt.Printf("❌ Unable to read blocks filename %s: %s\n", filename, err)
			return err
		}

		seenBlockCount := 0
		for {
			block, err := readerFactory.Read()
			if err != nil {
				if err == io.EOF {
					fmt.Fprintf(os.Stderr, "Total blocks: %d\n", seenBlockCount)
					return nil
				}
				return fmt.Errorf("error receiving blocks: %w", err)
			}

			seenBlockCount++

			if err := displayBlock(block, chain, outputPrinter); err != nil {
				// Error is ready to be passed to the user as-is
				return err
			}
		}
	}
}

func createToolsPrintOneBlockE[B firecore.Block](chain *firecore.Chain[B]) firecore.CommandExecutor {
	return func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		outputPrinter, err := GetOutputPrinter(cmd, chain.BlockFileDescriptor())
		cli.NoError(err, "Unable to get output printer")

		storeURL := args[0]
		store, err := dstore.NewDBinStore(storeURL)
		if err != nil {
			return fmt.Errorf("unable to create store at path %q: %w", store, err)
		}

		blockNum, err := strconv.ParseUint(args[1], 10, 64)
		if err != nil {
			return fmt.Errorf("unable to parse block number %q: %w", args[1], err)
		}

		var files []string
		filePrefix := fmt.Sprintf("%010d", blockNum)
		err = store.Walk(ctx, filePrefix, func(filename string) (err error) {
			files = append(files, filename)
			return nil
		})
		if err != nil {
			return fmt.Errorf("unable to find on block files: %w", err)
		}

		for _, filepath := range files {
			reader, err := store.OpenObject(ctx, filepath)
			if err != nil {
				fmt.Printf("❌ Unable to read block filename %s: %s\n", filepath, err)
				return err
			}
			defer reader.Close()

			readerFactory, err := bstream.NewDBinBlockReader(reader)
			if err != nil {
				fmt.Printf("❌ Unable to read blocks filename %s: %s\n", filepath, err)
				return err
			}

			block, err := readerFactory.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				return fmt.Errorf("reading block: %w", err)
			}

			if err := displayBlock(block, chain, outputPrinter); err != nil {
				// Error is ready to be passed to the user as-is
				return err
			}
		}
		return nil
	}
}

func displayBlock[B firecore.Block](pbBlock *pbbstream.Block, chain *firecore.Chain[B], printer OutputPrinter) error {
	if pbBlock == nil {
		return fmt.Errorf("block is nil")
	}

	if !firecore.UnsafeRunningFromFirecore {
		// since we are running via the chain specific binary (i.e. fireeth) we can use a BlockFactory
		marshallableBlock := chain.BlockFactory()

		if err := pbBlock.Payload.UnmarshalTo(marshallableBlock); err != nil {
			return fmt.Errorf("pbBlock payload unmarshal: %w", err)
		}

		err := printer.PrintTo(marshallableBlock, os.Stdout)
		if err != nil {
			return fmt.Errorf("pbBlock JSON printing: json marshal: %w", err)
		}
		return nil
	}

	// since we are running directly the firecore binary we will *NOT* use the BlockFactory
	err := printer.PrintTo(pbBlock.Payload, os.Stdout)
	if err != nil {
		return fmt.Errorf("marshalling block to json: %w", err)
	}

	return nil
}
