package node

//go:generate ruby ../prism/templates/template.rb ../../templates/gen_node.go ../parser/gen_node.go
//go:generate ruby ../prism/templates/template.rb ../../templates/gen_error.go ../parser/gen_error.go
//go:generate ruby ../prism/templates/template.rb ../../templates/gen_warning.go ../parser/gen_warning.go
//go:generate ruby ../prism/templates/template.rb ../../templates/gen_node_loader.go ../parser/gen_node_loader_node.go
//go:generate ruby ../prism/templates/template.rb ../../templates/gen_flags.go ../parser/gen_flags.go
//go:generate ruby ../prism/templates/template.rb ../../templates/gen_abstract_node_visitor.go ../parser/gen_abstract_node_visitor.go
