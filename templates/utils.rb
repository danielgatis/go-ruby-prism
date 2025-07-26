# frozen_string_literal: true

class String
  def capitalize
    self[0].upcase + self[1..]
  end

  def snake_case
    gsub(/([a-z])([A-Z])/, '\1_\2').downcase
  end

  def pascal_case
    split('_').map(&:capitalize).join
  end

  def camel_case
    split('_').map.with_index { |part, i| i.zero? ? part : part.capitalize }.join
  end
end

Prism::Template::NodeKindField.class_eval do
  def go_type
    if specific_kind
      "*#{specific_kind}"
    else
      'Node'
    end
  end
end

Prism::Template::ConstantField.class_eval do
  def go_type
    'string'
  end
end

Prism::Template::OptionalConstantField.class_eval do
  def go_type
    '*string'
  end
end

Prism::Template::ConstantListField.class_eval do
  def go_type
    '[]string'
  end
end

Prism::Template::DoubleField.class_eval do
  def go_type
    'float64'
  end
end

Prism::Template::IntegerField.class_eval do
  def go_type
    '*big.Int'
  end
end

Prism::Template::UInt32Field.class_eval do
  def go_type
    'uint32'
  end
end

Prism::Template::UInt8Field.class_eval do
  def go_type
    'uint8'
  end
end

Prism::Template::StringField.class_eval do
  def go_type
    'string'
  end
end

Prism::Template::NodeListField.class_eval do
  def go_type
    if specific_kind
      "[]*#{specific_kind}"
    else
      '[]Node'
    end
  end
end
