module github.com/streamingfast/firehose-core

go 1.22.0

require (
	buf.build/gen/go/bufbuild/reflect/connectrpc/go v1.16.1-20240117202343-bf8f65e8876c.1
	buf.build/gen/go/bufbuild/reflect/protocolbuffers/go v1.33.0-20240117202343-bf8f65e8876c.1
	connectrpc.com/connect v1.16.1
	github.com/ShinyTrinkets/overseer v0.3.0
	github.com/dustin/go-humanize v1.0.1
	github.com/go-json-experiment/json v0.0.0-20231013223334-54c864be5b8d
	github.com/hashicorp/go-multierror v1.1.1
	github.com/iancoleman/strcase v0.3.0
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51
	github.com/mostynb/go-grpc-compression v1.1.17
	github.com/prometheus/client_golang v1.16.0
	github.com/spf13/cobra v1.7.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.15.0
	github.com/streamingfast/bstream v0.0.2-0.20241108153156-a5c6bc006f41
	github.com/streamingfast/cli v0.0.4-0.20241119021815-815afa473375
	github.com/streamingfast/dauth v0.0.0-20240222213226-519afc16cf84
	github.com/streamingfast/dbin v0.9.1-0.20231117225723-59790c798e2c
	github.com/streamingfast/derr v0.0.0-20230515163924-8570aaa43fe1
	github.com/streamingfast/dgrpc v0.0.0-20240423143010-f36784700c9a
	github.com/streamingfast/dhammer v0.0.0-20230125192823-c34bbd561bd4
	github.com/streamingfast/dmetering v0.0.0-20241101155221-489f5a9d9139
	github.com/streamingfast/dmetrics v0.0.0-20230919161904-206fa8ebd545
	github.com/streamingfast/dstore v0.1.1-0.20241011152904-9acd6205dc14
	github.com/streamingfast/jsonpb v0.0.0-20210811021341-3670f0aa02d0
	github.com/streamingfast/logging v0.0.0-20230608130331-f22c91403091
	github.com/streamingfast/payment-gateway v0.0.0-20240426151444-581e930c76e2
	github.com/streamingfast/pbgo v0.0.6-0.20240823134334-812f6a16c5cb
	github.com/streamingfast/snapshotter v0.0.0-20230316190750-5bcadfde44d0
	github.com/streamingfast/substreams v1.11.3
	github.com/stretchr/testify v1.9.0
	github.com/test-go/testify v1.1.4
	go.uber.org/multierr v1.10.0
	go.uber.org/zap v1.26.0
	golang.org/x/exp v0.0.0-20231006140011-7918f672742d
	google.golang.org/grpc v1.64.0
	google.golang.org/protobuf v1.34.2
)

require (
	cloud.google.com/go/auth v0.6.1 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.2 // indirect
	connectrpc.com/grpchealth v1.3.0 // indirect
	connectrpc.com/grpcreflect v1.2.0 // indirect
	connectrpc.com/otelconnect v0.7.0 // indirect
	github.com/alecthomas/participle v0.7.1 // indirect
	github.com/aymanbagabas/go-osc52/v2 v2.0.1 // indirect
	github.com/bobg/go-generics/v3 v3.4.0 // indirect
	github.com/bufbuild/protocompile v0.4.0 // indirect
	github.com/charmbracelet/lipgloss v1.0.0 // indirect
	github.com/charmbracelet/x/ansi v0.4.2 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/google/pprof v0.0.0-20221203041831-ce31453925ec // indirect
	github.com/google/s2a-go v0.1.7 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/muesli/termenv v0.15.3-0.20240618155329-98d742f6907a // indirect
	github.com/protocolbuffers/protoscope v0.0.0-20221109213918-8e7a6aafa2c9 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/sercand/kuberesolver/v5 v5.1.1 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.49.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240617180043-68d350f18fd4 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240624140628-dc46fd24d27d // indirect
)

require (
	cloud.google.com/go v0.115.0 // indirect
	cloud.google.com/go/compute/metadata v0.3.0 // indirect
	cloud.google.com/go/iam v1.1.8 // indirect
	cloud.google.com/go/monitoring v1.19.0 // indirect
	cloud.google.com/go/storage v1.42.0 // indirect
	cloud.google.com/go/trace v1.10.7 // indirect
	contrib.go.opencensus.io/exporter/stackdriver v0.13.10 // indirect
	contrib.go.opencensus.io/exporter/zipkin v0.1.1 // indirect
	github.com/Azure/azure-pipeline-go v0.2.3 // indirect
	github.com/Azure/azure-storage-blob-go v0.14.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp v0.32.3 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace v1.15.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/internal/resourcemapping v0.39.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/propagator v0.0.0-20221018185641-36f91511cfd7 // indirect
	github.com/KimMachineGun/automemlimit v0.2.4
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/RoaringBitmap/roaring v1.9.1 // indirect
	github.com/ShinyTrinkets/meta-logger v0.2.0 // indirect
	github.com/abourget/llerrgroup v0.2.0
	github.com/aws/aws-sdk-go v1.44.325 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bits-and-blooms/bitset v1.12.0 // indirect
	github.com/blendle/zapdriver v1.3.2-0.20200203083823-9200777f8a3d // indirect
	github.com/census-instrumentation/opencensus-proto v0.4.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/chzyer/readline v1.5.0 // indirect
	github.com/cilium/ebpf v0.4.0 // indirect
	github.com/cncf/xds/go v0.0.0-20240318125728-8a4994d93e50 // indirect
	github.com/containerd/cgroups v1.0.4 // indirect
	github.com/coreos/go-systemd/v22 v22.5.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/emicklei/go-restful/v3 v3.8.0 // indirect
	github.com/envoyproxy/go-control-plane v0.12.0 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.0.4 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.19.5 // indirect
	github.com/go-openapi/swag v0.21.1 // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/gnostic v0.5.7-v3refs // indirect
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.2 // indirect
	github.com/googleapis/gax-go/v2 v2.12.5 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jhump/protoreflect v1.14.0
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/josephburnett/jd v1.7.1
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/lithammer/dedent v1.1.0 // indirect
	github.com/logrusorgru/aurora v2.0.3+incompatible // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/manifoldco/promptui v0.9.0 // indirect
	github.com/mattn/go-ieproxy v0.0.1 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mr-tron/base58 v1.2.0
	github.com/mschoch/smat v0.2.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/opencontainers/runtime-spec v1.0.2 // indirect
	github.com/openzipkin/zipkin-go v0.4.2 // indirect
	github.com/paulbellamy/ratecounter v0.2.0 // indirect
	github.com/pelletier/go-toml/v2 v2.0.6 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.44.0 // indirect
	github.com/prometheus/procfs v0.11.0 // indirect
	github.com/rs/cors v1.10.0 // indirect
	github.com/schollz/closestmatch v2.1.0+incompatible // indirect
	github.com/sethvargo/go-retry v0.2.3 // indirect
	github.com/shopspring/decimal v1.3.1 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/spf13/afero v1.10.0 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/streamingfast/dtracing v0.0.0-20220305214756-b5c0e8699839
	github.com/streamingfast/opaque v0.0.0-20210811180740-0c01d37ea308 // indirect
	github.com/streamingfast/sf-tracing v0.0.0-20240430173521-888827872b90
	github.com/streamingfast/shutter v1.5.0
	github.com/subosito/gotenv v1.4.2 // indirect
	github.com/teris-io/shortid v0.0.0-20171029131806-771a37caa5cf // indirect
	github.com/tetratelabs/wazero v1.8.0 // indirect
	github.com/yourbasic/graph v0.0.0-20210606180040-8ecfec1c2869 // indirect
	go.opencensus.io v0.24.0
	go.opentelemetry.io/contrib/detectors/gcp v1.9.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.49.0
	go.opentelemetry.io/otel v1.24.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.23.1 // indirect
	go.opentelemetry.io/otel/exporters/zipkin v1.23.1 // indirect
	go.opentelemetry.io/otel/metric v1.24.0 // indirect
	go.opentelemetry.io/otel/sdk v1.24.0 // indirect
	go.opentelemetry.io/otel/trace v1.24.0 // indirect
	go.uber.org/atomic v1.10.0
	go.uber.org/automaxprocs v1.5.1
	golang.org/x/crypto v0.24.0 // indirect
	golang.org/x/mod v0.17.0 // indirect
	golang.org/x/net v0.26.0
	golang.org/x/oauth2 v0.21.0
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	golang.org/x/term v0.21.0 // indirect
	golang.org/x/text v0.18.0 // indirect
	golang.org/x/time v0.5.0 // indirect
	google.golang.org/api v0.187.0 // indirect
	google.golang.org/genproto v0.0.0-20240624140628-dc46fd24d27d // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/olivere/elastic.v3 v3.0.75
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/api v0.25.0 // indirect
	k8s.io/apimachinery v0.25.0 // indirect
	k8s.io/client-go v0.25.0 // indirect
	k8s.io/klog/v2 v2.70.1 // indirect
	k8s.io/kube-openapi v0.0.0-20220803162953-67bda5d908f1 // indirect
	k8s.io/utils v0.0.0-20220728103510-ee6ede2d64ed // indirect
	sigs.k8s.io/json v0.0.0-20220713155537-f223a00ba0e2 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)

replace (
	cloud.google.com/go => github.com/streamingfast/google-cloud-go v0.0.0-20241202194114-f77ff78d4f66
	github.com/ShinyTrinkets/overseer => github.com/streamingfast/overseer v0.2.1-0.20210326144022-ee491780e3ef
	github.com/bytecodealliance/wasmtime-go/v4 => github.com/streamingfast/wasmtime-go/v4 v4.0.0-freemem3
	github.com/jhump/protoreflect => github.com/streamingfast/protoreflect v0.0.0-20231205191344-4b629d20ce8d
	github.com/tetratelabs/wazero => github.com/streamingfast/wazero v0.0.0-20241202185309-91287c3640ed

)
