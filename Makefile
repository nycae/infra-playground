PROTOC := protoc
GOC := go build

PROTODIR := api/

PROTOFLG := --go_out=$(PROTODIR) --go_opt=paths=source_relative
GRPCFLG := --go-grpc_out=$(PROTODIR) --go-grpc_opt=paths=source_relative

all: protos

protos: people_grpc.pb.go people.pb.go

%_grpc.pb.go: $(PROTODIR)%.proto
	$(PROTOC) -I$(PROTODIR) $(GRPCFLG) $<

%.pb.go: $(PROTODIR)%.proto
	$(PROTOC) -I$(PROTODIR) $(PROTOFLG) $<

clean: clean-protos

clean-protos:
	rm $(PROTODIR)*.go
