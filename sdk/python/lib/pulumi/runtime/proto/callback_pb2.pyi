"""
@generated by mypy-protobuf.  Do not edit manually!
isort:skip_file
Copyright 2016-2023, Pulumi Corporation.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""
import builtins
import collections.abc
import google.protobuf.descriptor
import google.protobuf.internal.containers
import google.protobuf.message
import google.protobuf.struct_pb2
import sys

if sys.version_info >= (3, 8):
    import typing as typing_extensions
else:
    import typing_extensions

DESCRIPTOR: google.protobuf.descriptor.FileDescriptor

@typing_extensions.final
class Callback(google.protobuf.message.Message):
    DESCRIPTOR: google.protobuf.descriptor.Descriptor

    TARGET_FIELD_NUMBER: builtins.int
    TOKEN_FIELD_NUMBER: builtins.int
    target: builtins.str
    """the gRPC target of the callback service."""
    token: builtins.str
    """the service specific unique token for this callback."""
    def __init__(
        self,
        *,
        target: builtins.str = ...,
        token: builtins.str = ...,
    ) -> None: ...
    def ClearField(self, field_name: typing_extensions.Literal["target", b"target", "token", b"token"]) -> None: ...

global___Callback = Callback

@typing_extensions.final
class CallbackInvokeRequest(google.protobuf.message.Message):
    DESCRIPTOR: google.protobuf.descriptor.Descriptor

    TOKEN_FIELD_NUMBER: builtins.int
    ARGUMENTS_FIELD_NUMBER: builtins.int
    token: builtins.str
    """the token for the callback."""
    @property
    def arguments(self) -> google.protobuf.internal.containers.RepeatedCompositeFieldContainer[google.protobuf.struct_pb2.Value]:
        """the property value arguments to pass to the callback function."""
    def __init__(
        self,
        *,
        token: builtins.str = ...,
        arguments: collections.abc.Iterable[google.protobuf.struct_pb2.Value] | None = ...,
    ) -> None: ...
    def ClearField(self, field_name: typing_extensions.Literal["arguments", b"arguments", "token", b"token"]) -> None: ...

global___CallbackInvokeRequest = CallbackInvokeRequest

@typing_extensions.final
class CallbackInvokeResponse(google.protobuf.message.Message):
    DESCRIPTOR: google.protobuf.descriptor.Descriptor

    RETURNS_FIELD_NUMBER: builtins.int
    @property
    def returns(self) -> google.protobuf.internal.containers.RepeatedCompositeFieldContainer[google.protobuf.struct_pb2.Value]:
        """The return values from the callback function."""
    def __init__(
        self,
        *,
        returns: collections.abc.Iterable[google.protobuf.struct_pb2.Value] | None = ...,
    ) -> None: ...
    def ClearField(self, field_name: typing_extensions.Literal["returns", b"returns"]) -> None: ...

global___CallbackInvokeResponse = CallbackInvokeResponse
