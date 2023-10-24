# coding=utf-8
# *** WARNING: this file was generated by test. ***
# *** Do not edit by hand unless you're certain you know what you are doing! ***

import copy
import warnings
import pulumi
import pulumi.runtime
from typing import Any, Callable, Mapping, Optional, Sequence, Union, overload
from ... import _utilities

__all__ = [
    'ENIConfigSpecArgs',
]

@pulumi.input_type
class ENIConfigSpecArgs:
    def __init__(__self__, *,
                 security_groups: Optional[pulumi.Input[Sequence[pulumi.Input[str]]]] = None,
                 subnet: Optional[pulumi.Input[str]] = None):
        ENIConfigSpecArgs._configure(
            lambda key, value: pulumi.set(__self__, key, value),
            security_groups=security_groups,
            subnet=subnet,
        )
    @staticmethod
    def _configure(
             _setter: Callable[[Any, Any], None],
             security_groups: Optional[pulumi.Input[Sequence[pulumi.Input[str]]]] = None,
             subnet: Optional[pulumi.Input[str]] = None,
             opts: Optional[pulumi.ResourceOptions]=None):
        if security_groups is not None:
            _setter("security_groups", security_groups)
        if subnet is not None:
            _setter("subnet", subnet)

    @property
    @pulumi.getter(name="securityGroups")
    def security_groups(self) -> Optional[pulumi.Input[Sequence[pulumi.Input[str]]]]:
        return pulumi.get(self, "security_groups")

    @security_groups.setter
    def security_groups(self, value: Optional[pulumi.Input[Sequence[pulumi.Input[str]]]]):
        pulumi.set(self, "security_groups", value)

    @property
    @pulumi.getter
    def subnet(self) -> Optional[pulumi.Input[str]]:
        return pulumi.get(self, "subnet")

    @subnet.setter
    def subnet(self, value: Optional[pulumi.Input[str]]):
        pulumi.set(self, "subnet", value)


