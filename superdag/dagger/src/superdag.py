import os
import importlib
import networkx as nx


class Dependency(object):
    def __init__(self, requirements_auto_discovery, artifact_requirements_type, **kwargs):
        """
            Args:
                requirements_auto_discovery (dict of list of string): represent a config where requirements module
                                                                      are bind on prefix path pattern
                artifact_requirements_type (dict of list of string): represent a config where artifact have specific
                                                                     requirements module set

            Keyword args:
                exclude_dependency_from_artifacts (set of string): represent a set of artifact to exclude from the graph
        """
        self.__kwargs = kwargs
        self.execution_by_level = list(list())  # list of dependencies level, top level to the lower level
        self.graph = None  # Represent a directed acyclic graph (networkx.classes.digraph.DiGraph)
        self.exclude_dependency_from_artifacts = kwargs.get("exclude_dependency_from_artifacts", set())
        self.requirements_auto_discovery = requirements_auto_discovery
        self.artifact_requirements_type = artifact_requirements_type

    def compute(self, root_path_to_analyze, modified_artifacts):
        """Compute modified artifacts in order to generate an order of artifact to execute base on a dag

        Args:
            root_path_to_analyze (str): the root path where the compute will begin
            modified_artifacts (set): artifact names that are considered modified

        Returns:
             list of list of string: represent each level of execution to follow
        """
        self._gen_dependencies_graph_from_modified_artifacts(root_path_to_analyze, modified_artifacts)

        self._graph_to_dependency_levels()

    def _graph_edges(self, artifacts_requirements, modified_artifacts):
        """Recursively find graph edges
        Args:
            artifacts_requirements (dict of str set): Dictionary of artifact_name: requirements
            modified_artifacts (set): artifact names that are considered modified
        Returns:
            list of tuple: Graph edges between artifacts.
        """
        if not modified_artifacts:
            return []

        edges_list = []
        next_modified_artifacts = set()

        for artifact_name in artifacts_requirements:
            requirements = artifacts_requirements[artifact_name]

            # but whose dependencies are not excluded (if this was not a modified artifact)
            modified_requirements = (requirements & modified_artifacts - self.exclude_dependency_from_artifacts
                                     if artifact_name not in modified_artifacts
                                     else requirements & modified_artifacts)
            if modified_requirements:
                # Add the edges from the modified requirements towards this artifact
                edges = map(lambda x: (x, artifact_name), modified_requirements)
                edges_list.extend(edges)
                next_modified_artifacts |= {artifact_name}

        edges_list.extend(self._graph_edges(artifacts_requirements, next_modified_artifacts))
        return edges_list

    def _get_artifact_requirements(self, artifact):
        """Fetch artifact requirements for an artifact

        Args:
            artifact (str): the artifact name / path

        Return:
            set of path dependencies each item is a path
        """
        deps = set()
        for cfg in self.requirements_auto_discovery:
            if artifact.startswith(cfg):
                for deps_module_type in self.requirements_auto_discovery[cfg]:
                    deps.update(
                        importlib.import_module(".dependencies.{}".
                                                format(deps_module_type)).RUN(artifact, **self.__kwargs)
                    )

        if artifact in self.artifact_requirements_type:
            for deps_module_type in self.artifact_requirements_type[artifact]:
                deps.update(
                    importlib.import_module(".dependencies.{}".
                                            format(deps_module_type)).RUN(artifact, **self.__kwargs)
                )

        return deps

    def _gen_dependencies_graph_from_modified_artifacts(self, root_path_to_analyze, modified_artifacts):
        """Parse requirement with a pattern for generating a dependencies graph
        Args:
            root_path_to_analyze (str): this is the path where the analysis has to begin
            modified_artifacts (set): This is a set a modified artifact
        """
        modified_artifacts = set(modified_artifacts)
        artifacts_requirements = dict()
        artifacts = ["{}/{}".format(root_path_to_analyze, stack)for stack in os.listdir(root_path_to_analyze)]

        # Go through all the artifacts
        for artifact in artifacts:
            requirements = self._get_artifact_requirements(artifact)
            if requirements:
                artifacts_requirements[artifact] = requirements

        edges_list = self._graph_edges(artifacts_requirements, modified_artifacts)

        directed_graph = nx.DiGraph(edges_list)
        directed_graph.add_nodes_from(modified_artifacts)

        self.graph = directed_graph

    def _graph_to_dependency_levels(self):
        """Generate relationship level dependencies with a directed graph
        """
        top_sort = list(nx.topological_sort(self.graph))

        levels = [[] for _ in range(len(top_sort))]

        for level, node in enumerate(top_sort):
            potential_level = level
            # While we are not connected to this level we push the artifact down a level
            while all((level_node, node) not in self.graph.edges() for level_node in levels[potential_level]) \
                    and potential_level >= 0:
                potential_level -= 1

            # Potential_level is now below the one where the node should go (where we are not connected)
            potential_level += 1
            levels[potential_level].append(node)

        self.execution_by_level = list(filter(None, levels))

    def filter(self, prefix):
        """Cleanup the output level in order to have only stacks in the execution plan

        Args:
            prefix (str): prefix to use for filtering in the result of dependencies

        Return:
            Return a list of dependencies level, top level to the lower level which match the prefix
        """
        return [
            keep_level_not_empty
            for keep_level_not_empty in [
                [artifact for artifact in level if artifact.startswith(prefix)]
                for level in self.execution_by_level
            ] if keep_level_not_empty
        ]
