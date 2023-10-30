package auth

// Code generated by Makefile; DO NOT EDIT.

var authModel = `{"schema_version":"1.1","type_definitions":[{"type":"user","relations":{}},{"type":"group","relations":{"member":{"this":{}}},"metadata":{"relations":{"member":{"directly_related_user_types":[{"type":"user"}]}}}},{"type":"server","relations":{"admin":{"this":{}},"operator":{"union":{"child":[{"this":{}},{"computedUserset":{"object":"","relation":"admin"}}]}},"viewer":{"union":{"child":[{"this":{}},{"computedUserset":{"object":"","relation":"operator"}}]}},"user":{"this":{}},"can_edit":{"computedUserset":{"object":"","relation":"admin"}},"can_view":{"computedUserset":{"object":"","relation":"user"}},"can_create_storage_pools":{"union":{"child":[{"this":{}},{"computedUserset":{"object":"","relation":"admin"}}]}},"can_create_projects":{"union":{"child":[{"this":{}},{"computedUserset":{"object":"","relation":"operator"}}]}},"can_view_resources":{"union":{"child":[{"this":{}},{"computedUserset":{"object":"","relation":"viewer"}}]}},"can_create_certificates":{"union":{"child":[{"this":{}},{"computedUserset":{"object":"","relation":"admin"}}]}},"can_view_metrics":{"union":{"child":[{"this":{}},{"computedUserset":{"object":"","relation":"viewer"}}]}},"can_override_cluster_target_restriction":{"union":{"child":[{"this":{}},{"computedUserset":{"object":"","relation":"admin"}}]}},"can_view_privileged_events":{"union":{"child":[{"this":{}},{"computedUserset":{"object":"","relation":"admin"}}]}}},"metadata":{"relations":{"admin":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"operator":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"viewer":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"user":{"directly_related_user_types":[{"type":"user","wildcard":{}}]},"can_edit":{"directly_related_user_types":[]},"can_view":{"directly_related_user_types":[]},"can_create_storage_pools":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"can_create_projects":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"can_view_resources":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"can_create_certificates":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"can_view_metrics":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"can_override_cluster_target_restriction":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"can_view_privileged_events":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]}}}},{"type":"certificate","relations":{"server":{"this":{}},"can_edit":{"union":{"child":[{"this":{}},{"tupleToUserset":{"tupleset":{"object":"","relation":"server"},"computedUserset":{"object":"","relation":"admin"}}}]}},"can_view":{"tupleToUserset":{"tupleset":{"object":"","relation":"server"},"computedUserset":{"object":"","relation":"user"}}}},"metadata":{"relations":{"server":{"directly_related_user_types":[{"type":"server"}]},"can_edit":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"can_view":{"directly_related_user_types":[]}}}},{"type":"storage_pool","relations":{"server":{"this":{}},"can_edit":{"union":{"child":[{"this":{}},{"tupleToUserset":{"tupleset":{"object":"","relation":"server"},"computedUserset":{"object":"","relation":"admin"}}}]}},"can_view":{"tupleToUserset":{"tupleset":{"object":"","relation":"server"},"computedUserset":{"object":"","relation":"user"}}}},"metadata":{"relations":{"server":{"directly_related_user_types":[{"type":"server"}]},"can_edit":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"can_view":{"directly_related_user_types":[]}}}},{"type":"project","relations":{"server":{"this":{}},"manager":{"union":{"child":[{"this":{}},{"tupleToUserset":{"tupleset":{"object":"","relation":"server"},"computedUserset":{"object":"","relation":"operator"}}}]}},"operator":{"union":{"child":[{"this":{}},{"computedUserset":{"object":"","relation":"manager"}},{"tupleToUserset":{"tupleset":{"object":"","relation":"server"},"computedUserset":{"object":"","relation":"operator"}}}]}},"viewer":{"union":{"child":[{"this":{}},{"computedUserset":{"object":"","relation":"operator"}}]}},"can_edit":{"computedUserset":{"object":"","relation":"manager"}},"can_view":{"computedUserset":{"object":"","relation":"viewer"}},"can_create_images":{"union":{"child":[{"this":{}},{"computedUserset":{"object":"","relation":"operator"}},{"tupleToUserset":{"tupleset":{"object":"","relation":"server"},"computedUserset":{"object":"","relation":"operator"}}}]}},"can_create_image_aliases":{"union":{"child":[{"this":{}},{"computedUserset":{"object":"","relation":"operator"}},{"tupleToUserset":{"tupleset":{"object":"","relation":"server"},"computedUserset":{"object":"","relation":"operator"}}}]}},"can_create_instances":{"union":{"child":[{"this":{}},{"computedUserset":{"object":"","relation":"operator"}},{"tupleToUserset":{"tupleset":{"object":"","relation":"server"},"computedUserset":{"object":"","relation":"operator"}}}]}},"can_create_networks":{"union":{"child":[{"this":{}},{"computedUserset":{"object":"","relation":"operator"}},{"tupleToUserset":{"tupleset":{"object":"","relation":"server"},"computedUserset":{"object":"","relation":"operator"}}}]}},"can_create_network_acls":{"union":{"child":[{"this":{}},{"computedUserset":{"object":"","relation":"operator"}},{"tupleToUserset":{"tupleset":{"object":"","relation":"server"},"computedUserset":{"object":"","relation":"operator"}}}]}},"can_create_network_zones":{"union":{"child":[{"this":{}},{"computedUserset":{"object":"","relation":"operator"}},{"tupleToUserset":{"tupleset":{"object":"","relation":"server"},"computedUserset":{"object":"","relation":"operator"}}}]}},"can_create_profiles":{"union":{"child":[{"this":{}},{"computedUserset":{"object":"","relation":"operator"}},{"tupleToUserset":{"tupleset":{"object":"","relation":"server"},"computedUserset":{"object":"","relation":"operator"}}}]}},"can_create_storage_volumes":{"union":{"child":[{"this":{}},{"computedUserset":{"object":"","relation":"operator"}},{"tupleToUserset":{"tupleset":{"object":"","relation":"server"},"computedUserset":{"object":"","relation":"operator"}}}]}},"can_create_storage_buckets":{"union":{"child":[{"this":{}},{"computedUserset":{"object":"","relation":"operator"}},{"tupleToUserset":{"tupleset":{"object":"","relation":"server"},"computedUserset":{"object":"","relation":"operator"}}}]}},"can_create_deployments":{"union":{"child":[{"this":{}},{"computedUserset":{"object":"","relation":"operator"}},{"tupleToUserset":{"tupleset":{"object":"","relation":"server"},"computedUserset":{"object":"","relation":"operator"}}}]}},"can_view_operations":{"union":{"child":[{"this":{}},{"computedUserset":{"object":"","relation":"viewer"}}]}},"can_view_events":{"union":{"child":[{"this":{}},{"computedUserset":{"object":"","relation":"viewer"}}]}}},"metadata":{"relations":{"server":{"directly_related_user_types":[{"type":"server"}]},"manager":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"operator":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"viewer":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"can_edit":{"directly_related_user_types":[]},"can_view":{"directly_related_user_types":[]},"can_create_images":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"can_create_image_aliases":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"can_create_instances":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"can_create_networks":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"can_create_network_acls":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"can_create_network_zones":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"can_create_profiles":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"can_create_storage_volumes":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"can_create_storage_buckets":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"can_create_deployments":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"can_view_operations":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"can_view_events":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]}}}},{"type":"image","relations":{"project":{"this":{}},"can_edit":{"union":{"child":[{"this":{}},{"tupleToUserset":{"tupleset":{"object":"","relation":"project"},"computedUserset":{"object":"","relation":"operator"}}}]}},"can_view":{"union":{"child":[{"this":{}},{"computedUserset":{"object":"","relation":"can_edit"}},{"tupleToUserset":{"tupleset":{"object":"","relation":"project"},"computedUserset":{"object":"","relation":"viewer"}}}]}}},"metadata":{"relations":{"project":{"directly_related_user_types":[{"type":"project"}]},"can_edit":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"can_view":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]}}}},{"type":"image_alias","relations":{"project":{"this":{}},"can_edit":{"union":{"child":[{"this":{}},{"tupleToUserset":{"tupleset":{"object":"","relation":"project"},"computedUserset":{"object":"","relation":"operator"}}}]}},"can_view":{"union":{"child":[{"this":{}},{"computedUserset":{"object":"","relation":"can_edit"}},{"tupleToUserset":{"tupleset":{"object":"","relation":"project"},"computedUserset":{"object":"","relation":"viewer"}}}]}}},"metadata":{"relations":{"project":{"directly_related_user_types":[{"type":"project"}]},"can_edit":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"can_view":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]}}}},{"type":"instance","relations":{"project":{"this":{}},"manager":{"this":{}},"operator":{"union":{"child":[{"this":{}},{"computedUserset":{"object":"","relation":"manager"}}]}},"user":{"union":{"child":[{"this":{}},{"computedUserset":{"object":"","relation":"operator"}}]}},"viewer":{"union":{"child":[{"this":{}},{"computedUserset":{"object":"","relation":"operator"}}]}},"can_edit":{"union":{"child":[{"computedUserset":{"object":"","relation":"manager"}},{"tupleToUserset":{"tupleset":{"object":"","relation":"project"},"computedUserset":{"object":"","relation":"operator"}}}]}},"can_view":{"union":{"child":[{"computedUserset":{"object":"","relation":"user"}},{"computedUserset":{"object":"","relation":"viewer"}},{"tupleToUserset":{"tupleset":{"object":"","relation":"project"},"computedUserset":{"object":"","relation":"viewer"}}}]}},"can_update_state":{"union":{"child":[{"this":{}},{"computedUserset":{"object":"","relation":"operator"}},{"tupleToUserset":{"tupleset":{"object":"","relation":"project"},"computedUserset":{"object":"","relation":"operator"}}}]}},"can_manage_snapshots":{"union":{"child":[{"this":{}},{"computedUserset":{"object":"","relation":"operator"}},{"tupleToUserset":{"tupleset":{"object":"","relation":"project"},"computedUserset":{"object":"","relation":"operator"}}}]}},"can_manage_backups":{"union":{"child":[{"this":{}},{"computedUserset":{"object":"","relation":"operator"}},{"tupleToUserset":{"tupleset":{"object":"","relation":"project"},"computedUserset":{"object":"","relation":"operator"}}}]}},"can_connect_sftp":{"union":{"child":[{"this":{}},{"computedUserset":{"object":"","relation":"user"}},{"tupleToUserset":{"tupleset":{"object":"","relation":"project"},"computedUserset":{"object":"","relation":"operator"}}}]}},"can_access_files":{"union":{"child":[{"this":{}},{"computedUserset":{"object":"","relation":"user"}},{"tupleToUserset":{"tupleset":{"object":"","relation":"project"},"computedUserset":{"object":"","relation":"operator"}}}]}},"can_access_console":{"union":{"child":[{"this":{}},{"computedUserset":{"object":"","relation":"user"}},{"tupleToUserset":{"tupleset":{"object":"","relation":"project"},"computedUserset":{"object":"","relation":"operator"}}}]}},"can_exec":{"union":{"child":[{"this":{}},{"computedUserset":{"object":"","relation":"user"}},{"tupleToUserset":{"tupleset":{"object":"","relation":"project"},"computedUserset":{"object":"","relation":"operator"}}}]}}},"metadata":{"relations":{"project":{"directly_related_user_types":[{"type":"project"}]},"manager":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"operator":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"user":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"viewer":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"can_edit":{"directly_related_user_types":[]},"can_view":{"directly_related_user_types":[]},"can_update_state":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"can_manage_snapshots":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"can_manage_backups":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"can_connect_sftp":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"can_access_files":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"can_access_console":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"can_exec":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]}}}},{"type":"network","relations":{"project":{"this":{}},"can_edit":{"union":{"child":[{"this":{}},{"tupleToUserset":{"tupleset":{"object":"","relation":"project"},"computedUserset":{"object":"","relation":"operator"}}}]}},"can_view":{"union":{"child":[{"this":{}},{"computedUserset":{"object":"","relation":"can_edit"}},{"tupleToUserset":{"tupleset":{"object":"","relation":"project"},"computedUserset":{"object":"","relation":"viewer"}}}]}}},"metadata":{"relations":{"project":{"directly_related_user_types":[{"type":"project"}]},"can_edit":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"can_view":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]}}}},{"type":"network_acl","relations":{"project":{"this":{}},"can_edit":{"union":{"child":[{"this":{}},{"tupleToUserset":{"tupleset":{"object":"","relation":"project"},"computedUserset":{"object":"","relation":"operator"}}}]}},"can_view":{"union":{"child":[{"this":{}},{"computedUserset":{"object":"","relation":"can_edit"}},{"tupleToUserset":{"tupleset":{"object":"","relation":"project"},"computedUserset":{"object":"","relation":"viewer"}}}]}}},"metadata":{"relations":{"project":{"directly_related_user_types":[{"type":"project"}]},"can_edit":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"can_view":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]}}}},{"type":"network_zone","relations":{"project":{"this":{}},"can_edit":{"union":{"child":[{"this":{}},{"tupleToUserset":{"tupleset":{"object":"","relation":"project"},"computedUserset":{"object":"","relation":"operator"}}}]}},"can_view":{"union":{"child":[{"this":{}},{"computedUserset":{"object":"","relation":"can_edit"}},{"tupleToUserset":{"tupleset":{"object":"","relation":"project"},"computedUserset":{"object":"","relation":"viewer"}}}]}}},"metadata":{"relations":{"project":{"directly_related_user_types":[{"type":"project"}]},"can_edit":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"can_view":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]}}}},{"type":"profile","relations":{"project":{"this":{}},"can_edit":{"union":{"child":[{"this":{}},{"tupleToUserset":{"tupleset":{"object":"","relation":"project"},"computedUserset":{"object":"","relation":"operator"}}}]}},"can_view":{"union":{"child":[{"this":{}},{"computedUserset":{"object":"","relation":"can_edit"}},{"tupleToUserset":{"tupleset":{"object":"","relation":"project"},"computedUserset":{"object":"","relation":"viewer"}}}]}}},"metadata":{"relations":{"project":{"directly_related_user_types":[{"type":"project"}]},"can_edit":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"can_view":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]}}}},{"type":"storage_volume","relations":{"project":{"this":{}},"can_edit":{"union":{"child":[{"this":{}},{"tupleToUserset":{"tupleset":{"object":"","relation":"project"},"computedUserset":{"object":"","relation":"operator"}}}]}},"can_view":{"union":{"child":[{"this":{}},{"computedUserset":{"object":"","relation":"can_edit"}},{"tupleToUserset":{"tupleset":{"object":"","relation":"project"},"computedUserset":{"object":"","relation":"viewer"}}}]}},"can_manage_snapshots":{"union":{"child":[{"this":{}},{"computedUserset":{"object":"","relation":"can_edit"}}]}},"can_manage_backups":{"union":{"child":[{"this":{}},{"computedUserset":{"object":"","relation":"can_edit"}}]}}},"metadata":{"relations":{"project":{"directly_related_user_types":[{"type":"project"}]},"can_edit":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"can_view":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"can_manage_snapshots":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"can_manage_backups":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]}}}},{"type":"storage_bucket","relations":{"project":{"this":{}},"can_edit":{"union":{"child":[{"this":{}},{"tupleToUserset":{"tupleset":{"object":"","relation":"project"},"computedUserset":{"object":"","relation":"operator"}}}]}},"can_view":{"union":{"child":[{"this":{}},{"computedUserset":{"object":"","relation":"can_edit"}},{"tupleToUserset":{"tupleset":{"object":"","relation":"project"},"computedUserset":{"object":"","relation":"viewer"}}}]}}},"metadata":{"relations":{"project":{"directly_related_user_types":[{"type":"project"}]},"can_edit":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"can_view":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]}}}},{"type":"deployment","relations":{"project":{"this":{}},"can_edit":{"union":{"child":[{"this":{}},{"tupleToUserset":{"tupleset":{"object":"","relation":"project"},"computedUserset":{"object":"","relation":"operator"}}}]}},"can_view":{"union":{"child":[{"this":{}},{"computedUserset":{"object":"","relation":"can_edit"}},{"tupleToUserset":{"tupleset":{"object":"","relation":"project"},"computedUserset":{"object":"","relation":"viewer"}}}]}},"can_access_deployment_keys":{"union":{"child":[{"this":{}},{"computedUserset":{"object":"","relation":"can_edit"}},{"tupleToUserset":{"tupleset":{"object":"","relation":"project"},"computedUserset":{"object":"","relation":"viewer"}}}]}},"can_create_deployment_keys":{"union":{"child":[{"this":{}},{"tupleToUserset":{"tupleset":{"object":"","relation":"project"},"computedUserset":{"object":"","relation":"operator"}}}]}},"can_access_deployment_shapes":{"union":{"child":[{"this":{}},{"computedUserset":{"object":"","relation":"can_edit"}},{"tupleToUserset":{"tupleset":{"object":"","relation":"project"},"computedUserset":{"object":"","relation":"viewer"}}}]}},"can_create_deployment_shapes":{"union":{"child":[{"this":{}},{"tupleToUserset":{"tupleset":{"object":"","relation":"project"},"computedUserset":{"object":"","relation":"operator"}}}]}}},"metadata":{"relations":{"project":{"directly_related_user_types":[{"type":"project"}]},"can_edit":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"can_view":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"can_access_deployment_keys":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"can_create_deployment_keys":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"can_access_deployment_shapes":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"can_create_deployment_shapes":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]}}}},{"type":"deployment_shape","relations":{"project":{"this":{}},"can_edit":{"union":{"child":[{"this":{}},{"tupleToUserset":{"tupleset":{"object":"","relation":"project"},"computedUserset":{"object":"","relation":"operator"}}}]}},"can_view":{"union":{"child":[{"this":{}},{"computedUserset":{"object":"","relation":"can_edit"}},{"tupleToUserset":{"tupleset":{"object":"","relation":"project"},"computedUserset":{"object":"","relation":"viewer"}}}]}},"can_access_deployed_instances":{"union":{"child":[{"this":{}},{"tupleToUserset":{"tupleset":{"object":"","relation":"project"},"computedUserset":{"object":"","relation":"viewer"}}}]}},"can_deploy_instances":{"union":{"child":[{"this":{}},{"computedUserset":{"object":"","relation":"can_edit"}},{"tupleToUserset":{"tupleset":{"object":"","relation":"project"},"computedUserset":{"object":"","relation":"viewer"}}}]}}},"metadata":{"relations":{"project":{"directly_related_user_types":[{"type":"project"}]},"can_edit":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"can_view":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"can_access_deployed_instances":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"can_deploy_instances":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]}}}},{"type":"deployment_key","relations":{"project":{"this":{}},"can_edit":{"union":{"child":[{"this":{}},{"tupleToUserset":{"tupleset":{"object":"","relation":"project"},"computedUserset":{"object":"","relation":"operator"}}}]}},"can_view":{"union":{"child":[{"this":{}},{"computedUserset":{"object":"","relation":"can_edit"}},{"tupleToUserset":{"tupleset":{"object":"","relation":"project"},"computedUserset":{"object":"","relation":"viewer"}}}]}}},"metadata":{"relations":{"project":{"directly_related_user_types":[{"type":"project"}]},"can_edit":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"can_view":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]}}}},{"type":"deployment_shape_instance","relations":{"project":{"this":{}},"can_edit":{"union":{"child":[{"this":{}},{"tupleToUserset":{"tupleset":{"object":"","relation":"project"},"computedUserset":{"object":"","relation":"operator"}}}]}},"can_view":{"union":{"child":[{"this":{}},{"computedUserset":{"object":"","relation":"can_edit"}},{"tupleToUserset":{"tupleset":{"object":"","relation":"project"},"computedUserset":{"object":"","relation":"viewer"}}}]}}},"metadata":{"relations":{"project":{"directly_related_user_types":[{"type":"project"}]},"can_edit":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]},"can_view":{"directly_related_user_types":[{"type":"user"},{"type":"group","relation":"member"}]}}}}]}`
