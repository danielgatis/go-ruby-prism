def const_prefix(name)
  name.gsub("Flags", "").gsub(/([a-z])([A-Z])/, '\1_\2').upcase
end

def arg(field)
  field.name.gsub(/_([a-z])/) { $1.upcase }
end

def prop(field)
  arg(field).capitalize
end

def gotype(field)
  case field
  when Prism::Template::NodeField then field.ruby_type == "Node" ? "Node" : "*#{field.ruby_type}"
  when Prism::Template::OptionalNodeField then field.ruby_type == "Node" ? "Node" : "*#{field.ruby_type}"
  when Prism::Template::NodeListField then "[]Node"
  when Prism::Template::StringField then "string"
  when Prism::Template::ConstantField then "string"
  when Prism::Template::OptionalConstantField then "*string"
  when Prism::Template::ConstantListField then "[]string"
  when Prism::Template::LocationField then "*Location"
  when Prism::Template::OptionalLocationField then "*Location"
  when Prism::Template::FlagsField then field.options[:kind]
  when Prism::Template::DoubleField then "float64"
  when Prism::Template::UInt8Field then "uint8"
  when Prism::Template::UInt32Field then "uint32"
  when Prism::Template::IntegerField then "*big.Int"
  else raise "Unknown field type: #{field.inspect}"
  end
end
