package node

//go:generate ruby ./template.rb gen_node.go ../parser/gen_node.go
//go:generate ruby ./template.rb gen_syntaxe_error.go ../parser/gen_syntaxe_error.go
//go:generate ruby ./template.rb gen_syntaxe_warn.go ../parser/gen_syntaxe_warn.go
//go:generate ruby ./template.rb gen_loader_node.go ../parser/gen_loader_node.go
//go:generate ruby ./template.rb gen_flags.go ../parser/gen_flags.go
