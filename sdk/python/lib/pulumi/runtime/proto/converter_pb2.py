# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: pulumi/converter.proto
"""Generated protocol buffer code."""
from google.protobuf.internal import builder as _builder
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from .codegen import hcl_pb2 as pulumi_dot_codegen_dot_hcl__pb2
from . import plugin_pb2 as pulumi_dot_plugin__pb2
from google.protobuf import empty_pb2 as google_dot_protobuf_dot_empty__pb2
from google.protobuf import struct_pb2 as google_dot_protobuf_dot_struct__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x16pulumi/converter.proto\x12\tpulumirpc\x1a\x18pulumi/codegen/hcl.proto\x1a\x13pulumi/plugin.proto\x1a\x1bgoogle/protobuf/empty.proto\x1a\x1cgoogle/protobuf/struct.proto\":\n\x13\x43onvertStateRequest\x12\x15\n\rmapper_target\x18\x01 \x01(\t\x12\x0c\n\x04\x61rgs\x18\x02 \x03(\t\"d\n\x0eResourceImport\x12\x0c\n\x04type\x18\x01 \x01(\t\x12\x0c\n\x04name\x18\x02 \x01(\t\x12\n\n\x02id\x18\x03 \x01(\t\x12\x0f\n\x07version\x18\x04 \x01(\t\x12\x19\n\x11pluginDownloadURL\x18\x05 \x01(\t\"D\n\x14\x43onvertStateResponse\x12,\n\tresources\x18\x01 \x03(\x0b\x32\x19.pulumirpc.ResourceImport\"\x87\x01\n\x15\x43onvertProgramRequest\x12\x18\n\x10source_directory\x18\x01 \x01(\t\x12\x18\n\x10target_directory\x18\x02 \x01(\t\x12\x15\n\rmapper_target\x18\x03 \x01(\t\x12\x15\n\rloader_target\x18\x04 \x01(\t\x12\x0c\n\x04\x61rgs\x18\x05 \x03(\t\"L\n\x16\x43onvertProgramResponse\x12\x32\n\x0b\x64iagnostics\x18\x01 \x03(\x0b\x32\x1d.pulumirpc.codegen.Diagnostic2\xb7\x01\n\tConverter\x12Q\n\x0c\x43onvertState\x12\x1e.pulumirpc.ConvertStateRequest\x1a\x1f.pulumirpc.ConvertStateResponse\"\x00\x12W\n\x0e\x43onvertProgram\x12 .pulumirpc.ConvertProgramRequest\x1a!.pulumirpc.ConvertProgramResponse\"\x00\x42\x34Z2github.com/pulumi/pulumi/sdk/v3/proto/go;pulumirpcb\x06proto3')

_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, globals())
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'pulumi.converter_pb2', globals())
if _descriptor._USE_C_DESCRIPTORS == False:

  DESCRIPTOR._options = None
  DESCRIPTOR._serialized_options = b'Z2github.com/pulumi/pulumi/sdk/v3/proto/go;pulumirpc'
  _CONVERTSTATEREQUEST._serialized_start=143
  _CONVERTSTATEREQUEST._serialized_end=201
  _RESOURCEIMPORT._serialized_start=203
  _RESOURCEIMPORT._serialized_end=303
  _CONVERTSTATERESPONSE._serialized_start=305
  _CONVERTSTATERESPONSE._serialized_end=373
  _CONVERTPROGRAMREQUEST._serialized_start=376
  _CONVERTPROGRAMREQUEST._serialized_end=511
  _CONVERTPROGRAMRESPONSE._serialized_start=513
  _CONVERTPROGRAMRESPONSE._serialized_end=589
  _CONVERTER._serialized_start=592
  _CONVERTER._serialized_end=775
# @@protoc_insertion_point(module_scope)
