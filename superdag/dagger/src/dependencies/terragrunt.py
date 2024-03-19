import hcl2
from itertools import chain
import os
import re




def extract_terragrunt(stack_root_path, stack_path, **kwargs):
    """Get terragrunt source code and dependencies from terragrunt.hcl.
    By default, it will import also terraform module, and terraform sops secret

    Args:
        stack_root_path(str): The common root path for all stacks
        stack_path (str): The sub path to the stack

    Keyword args:
        skip_terraform_deps (bool): indicate to not fetch terraform deps
        skip_terraform_sops_deps (bool): indicate to not fetch sops terraform deps
        pattern_to_replace(dict): contain pattern as key and the replacement as a value of the dict

    Returns:
        set: dependencies path.
    """

    placeholder_patterns = kwargs.get("pattern_to_replace", {})
    dep_names = set()
    terragrunt_manifest_path = os.path.join(stack_root_path+"/"+stack_path, "terragrunt.hcl")
    stack_information = list(filter(None, stack_path.split("/")))
    account = stack_information[3]
    region = stack_information[4]
    env = stack_information[5]

    if not os.path.exists(terragrunt_manifest_path):
        return set()

    with open(terragrunt_manifest_path, "r") as f:
        terragrunt_manifest_content = hcl2.load(f)

    terraform_source = terragrunt_manifest_content["source"].replace("//", "/")

    for k, v in placeholder_patterns:
        if k in terraform_source:
            terraform_source = terraform_source.replace(k, v)

    dep_names.add(terraform_source)
    stack_dependencies = list(chain.from_iterable([j for j in [i["paths"] for i in terragrunt_manifest_content["dependencies"]]]))
    dep_names.update(stack_dependencies)

    if not kwargs.get("skip_terraform_deps", False):
        dep_names.update(extract_terraform_modules(terraform_source))

    if not kwargs.get("skip_terraform_sops_deps", False):
        dep_names.update(extract_sops_from_terraform(terraform_source, account, region, env))

    return dep_names


RUN = extract_terragrunt