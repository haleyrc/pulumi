# coding=utf-8
# *** WARNING: this file was generated by test. ***
# *** Do not edit by hand unless you're certain you know what you are doing! ***

import copy
import warnings
import pulumi
import pulumi.runtime
from typing import Any, Callable, Mapping, Optional, Sequence, Union, overload
from . import _utilities

__all__ = [
    'Config',
]

@pulumi.output_type
class Config(dict):
    def __init__(__self__, *,
                 foo: Optional[str] = None):
        Config._configure(
            lambda key, value: pulumi.set(__self__, key, value),
            foo=foo,
        )
    @staticmethod
    def _configure(
             _setter: Callable[[Any, Any], None],
             foo: Optional[str] = None,
             opts: Optional[pulumi.ResourceOptions] = None,
             **kwargs):

        if foo is not None:
            _setter("foo", foo)

    @property
    @pulumi.getter
    def foo(self) -> Optional[str]:
        return pulumi.get(self, "foo")


