/******************************************************************************/
/* This file is generated by the templates/template.rb script and should not  */
/* be modified manually. See                                                  */
/* templates/gen_syntaxe_error.go.erb                                         */
/* if you are looking to modify the                                           */
/* template                                                                   */
/******************************************************************************/

package parser

type SyntaxErrorLevel int
type SyntaxErrorType int

const (
	SyntaxErrorFatal SyntaxErrorLevel = iota
	SyntaxErrorArgument
)

const (
	ALIAS_ARGUMENT                     SyntaxErrorType = 0
	AMPAMPEQ_MULTI_ASSIGN              SyntaxErrorType = 1
	ARGUMENT_AFTER_BLOCK               SyntaxErrorType = 2
	ARGUMENT_AFTER_FORWARDING_ELLIPSES SyntaxErrorType = 3
	ARGUMENT_BARE_HASH                 SyntaxErrorType = 4
	ARGUMENT_BLOCK_FORWARDING          SyntaxErrorType = 5
	ARGUMENT_BLOCK_MULTI               SyntaxErrorType = 6
	ARGUMENT_FORMAL_CLASS              SyntaxErrorType = 7
	ARGUMENT_FORMAL_CONSTANT           SyntaxErrorType = 8
	ARGUMENT_FORMAL_GLOBAL             SyntaxErrorType = 9
	ARGUMENT_FORMAL_IVAR               SyntaxErrorType = 10
	ARGUMENT_FORWARDING_UNBOUND        SyntaxErrorType = 11
	ARGUMENT_IN                        SyntaxErrorType = 12
	ARGUMENT_NO_FORWARDING_AMP         SyntaxErrorType = 13
	ARGUMENT_NO_FORWARDING_ELLIPSES    SyntaxErrorType = 14
	ARGUMENT_NO_FORWARDING_STAR        SyntaxErrorType = 15
	ARGUMENT_SPLAT_AFTER_ASSOC_SPLAT   SyntaxErrorType = 16
	ARGUMENT_SPLAT_AFTER_SPLAT         SyntaxErrorType = 17
	ARGUMENT_TERM_PAREN                SyntaxErrorType = 18
	ARGUMENT_UNEXPECTED_BLOCK          SyntaxErrorType = 19
	ARRAY_ELEMENT                      SyntaxErrorType = 20
	ARRAY_EXPRESSION                   SyntaxErrorType = 21
	ARRAY_EXPRESSION_AFTER_STAR        SyntaxErrorType = 22
	ARRAY_SEPARATOR                    SyntaxErrorType = 23
	ARRAY_TERM                         SyntaxErrorType = 24
	BEGIN_LONELY_ELSE                  SyntaxErrorType = 25
	BEGIN_TERM                         SyntaxErrorType = 26
	BEGIN_UPCASE_BRACE                 SyntaxErrorType = 27
	BEGIN_UPCASE_TERM                  SyntaxErrorType = 28
	BEGIN_UPCASE_TOPLEVEL              SyntaxErrorType = 29
	BLOCK_PARAM_LOCAL_VARIABLE         SyntaxErrorType = 30
	BLOCK_PARAM_PIPE_TERM              SyntaxErrorType = 31
	BLOCK_TERM_BRACE                   SyntaxErrorType = 32
	BLOCK_TERM_END                     SyntaxErrorType = 33
	CANNOT_PARSE_EXPRESSION            SyntaxErrorType = 34
	CANNOT_PARSE_STRING_PART           SyntaxErrorType = 35
	CASE_EXPRESSION_AFTER_CASE         SyntaxErrorType = 36
	CASE_EXPRESSION_AFTER_WHEN         SyntaxErrorType = 37
	CASE_MATCH_MISSING_PREDICATE       SyntaxErrorType = 38
	CASE_MISSING_CONDITIONS            SyntaxErrorType = 39
	CASE_TERM                          SyntaxErrorType = 40
	CLASS_IN_METHOD                    SyntaxErrorType = 41
	CLASS_NAME                         SyntaxErrorType = 42
	CLASS_SUPERCLASS                   SyntaxErrorType = 43
	CLASS_TERM                         SyntaxErrorType = 44
	CLASS_UNEXPECTED_END               SyntaxErrorType = 45
	CONDITIONAL_ELSIF_PREDICATE        SyntaxErrorType = 46
	CONDITIONAL_IF_PREDICATE           SyntaxErrorType = 47
	CONDITIONAL_PREDICATE_TERM         SyntaxErrorType = 48
	CONDITIONAL_TERM                   SyntaxErrorType = 49
	CONDITIONAL_TERM_ELSE              SyntaxErrorType = 50
	CONDITIONAL_UNLESS_PREDICATE       SyntaxErrorType = 51
	CONDITIONAL_UNTIL_PREDICATE        SyntaxErrorType = 52
	CONDITIONAL_WHILE_PREDICATE        SyntaxErrorType = 53
	CONSTANT_PATH_COLON_COLON_CONSTANT SyntaxErrorType = 54
	DEF_ENDLESS                        SyntaxErrorType = 55
	DEF_ENDLESS_SETTER                 SyntaxErrorType = 56
	DEF_NAME                           SyntaxErrorType = 57
	DEF_NAME_AFTER_RECEIVER            SyntaxErrorType = 58
	DEF_PARAMS_TERM                    SyntaxErrorType = 59
	DEF_PARAMS_TERM_PAREN              SyntaxErrorType = 60
	DEF_RECEIVER                       SyntaxErrorType = 61
	DEF_RECEIVER_TERM                  SyntaxErrorType = 62
	DEF_TERM                           SyntaxErrorType = 63
	DEFINED_EXPRESSION                 SyntaxErrorType = 64
	EMBDOC_TERM                        SyntaxErrorType = 65
	EMBEXPR_END                        SyntaxErrorType = 66
	EMBVAR_INVALID                     SyntaxErrorType = 67
	END_UPCASE_BRACE                   SyntaxErrorType = 68
	END_UPCASE_TERM                    SyntaxErrorType = 69
	ESCAPE_INVALID_CONTROL             SyntaxErrorType = 70
	ESCAPE_INVALID_CONTROL_REPEAT      SyntaxErrorType = 71
	ESCAPE_INVALID_HEXADECIMAL         SyntaxErrorType = 72
	ESCAPE_INVALID_META                SyntaxErrorType = 73
	ESCAPE_INVALID_META_REPEAT         SyntaxErrorType = 74
	ESCAPE_INVALID_UNICODE             SyntaxErrorType = 75
	ESCAPE_INVALID_UNICODE_CM_FLAGS    SyntaxErrorType = 76
	ESCAPE_INVALID_UNICODE_LITERAL     SyntaxErrorType = 77
	ESCAPE_INVALID_UNICODE_LONG        SyntaxErrorType = 78
	ESCAPE_INVALID_UNICODE_TERM        SyntaxErrorType = 79
	EXPECT_ARGUMENT                    SyntaxErrorType = 80
	EXPECT_EOL_AFTER_STATEMENT         SyntaxErrorType = 81
	EXPECT_EXPRESSION_AFTER_AMPAMPEQ   SyntaxErrorType = 82
	EXPECT_EXPRESSION_AFTER_COMMA      SyntaxErrorType = 83
	EXPECT_EXPRESSION_AFTER_EQUAL      SyntaxErrorType = 84
	EXPECT_EXPRESSION_AFTER_LESS_LESS  SyntaxErrorType = 85
	EXPECT_EXPRESSION_AFTER_LPAREN     SyntaxErrorType = 86
	EXPECT_EXPRESSION_AFTER_OPERATOR   SyntaxErrorType = 87
	EXPECT_EXPRESSION_AFTER_PIPEPIPEEQ SyntaxErrorType = 88
	EXPECT_EXPRESSION_AFTER_QUESTION   SyntaxErrorType = 89
	EXPECT_EXPRESSION_AFTER_SPLAT      SyntaxErrorType = 90
	EXPECT_EXPRESSION_AFTER_SPLAT_HASH SyntaxErrorType = 91
	EXPECT_EXPRESSION_AFTER_STAR       SyntaxErrorType = 92
	EXPECT_IDENT_REQ_PARAMETER         SyntaxErrorType = 93
	EXPECT_LPAREN_REQ_PARAMETER        SyntaxErrorType = 94
	EXPECT_RBRACKET                    SyntaxErrorType = 95
	EXPECT_RPAREN                      SyntaxErrorType = 96
	EXPECT_RPAREN_AFTER_MULTI          SyntaxErrorType = 97
	EXPECT_RPAREN_REQ_PARAMETER        SyntaxErrorType = 98
	EXPECT_STRING_CONTENT              SyntaxErrorType = 99
	EXPECT_WHEN_DELIMITER              SyntaxErrorType = 100
	EXPRESSION_BARE_HASH               SyntaxErrorType = 101
	FLOAT_PARSE                        SyntaxErrorType = 102
	FOR_COLLECTION                     SyntaxErrorType = 103
	FOR_IN                             SyntaxErrorType = 104
	FOR_INDEX                          SyntaxErrorType = 105
	FOR_TERM                           SyntaxErrorType = 106
	HASH_EXPRESSION_AFTER_LABEL        SyntaxErrorType = 107
	HASH_KEY                           SyntaxErrorType = 108
	HASH_ROCKET                        SyntaxErrorType = 109
	HASH_TERM                          SyntaxErrorType = 110
	HASH_VALUE                         SyntaxErrorType = 111
	HEREDOC_TERM                       SyntaxErrorType = 112
	INCOMPLETE_QUESTION_MARK           SyntaxErrorType = 113
	INCOMPLETE_VARIABLE_CLASS          SyntaxErrorType = 114
	INCOMPLETE_VARIABLE_CLASS_3_3_0    SyntaxErrorType = 115
	INCOMPLETE_VARIABLE_INSTANCE       SyntaxErrorType = 116
	INCOMPLETE_VARIABLE_INSTANCE_3_3_0 SyntaxErrorType = 117
	INVALID_CHARACTER                  SyntaxErrorType = 118
	INVALID_ENCODING_MAGIC_COMMENT     SyntaxErrorType = 119
	INVALID_FLOAT_EXPONENT             SyntaxErrorType = 120
	INVALID_MULTIBYTE_CHAR             SyntaxErrorType = 121
	INVALID_MULTIBYTE_CHARACTER        SyntaxErrorType = 122
	INVALID_MULTIBYTE_ESCAPE           SyntaxErrorType = 123
	INVALID_NUMBER_BINARY              SyntaxErrorType = 124
	INVALID_NUMBER_DECIMAL             SyntaxErrorType = 125
	INVALID_NUMBER_HEXADECIMAL         SyntaxErrorType = 126
	INVALID_NUMBER_OCTAL               SyntaxErrorType = 127
	INVALID_NUMBER_UNDERSCORE          SyntaxErrorType = 128
	INVALID_PERCENT                    SyntaxErrorType = 129
	INVALID_PRINTABLE_CHARACTER        SyntaxErrorType = 130
	INVALID_VARIABLE_GLOBAL            SyntaxErrorType = 131
	INVALID_VARIABLE_GLOBAL_3_3_0      SyntaxErrorType = 132
	IT_NOT_ALLOWED_NUMBERED            SyntaxErrorType = 133
	IT_NOT_ALLOWED_ORDINARY            SyntaxErrorType = 134
	LAMBDA_OPEN                        SyntaxErrorType = 135
	LAMBDA_TERM_BRACE                  SyntaxErrorType = 136
	LAMBDA_TERM_END                    SyntaxErrorType = 137
	LIST_I_LOWER_ELEMENT               SyntaxErrorType = 138
	LIST_I_LOWER_TERM                  SyntaxErrorType = 139
	LIST_I_UPPER_ELEMENT               SyntaxErrorType = 140
	LIST_I_UPPER_TERM                  SyntaxErrorType = 141
	LIST_W_LOWER_ELEMENT               SyntaxErrorType = 142
	LIST_W_LOWER_TERM                  SyntaxErrorType = 143
	LIST_W_UPPER_ELEMENT               SyntaxErrorType = 144
	LIST_W_UPPER_TERM                  SyntaxErrorType = 145
	MALLOC_FAILED                      SyntaxErrorType = 146
	MIXED_ENCODING                     SyntaxErrorType = 147
	MODULE_IN_METHOD                   SyntaxErrorType = 148
	MODULE_NAME                        SyntaxErrorType = 149
	MODULE_TERM                        SyntaxErrorType = 150
	MULTI_ASSIGN_MULTI_SPLATS          SyntaxErrorType = 151
	MULTI_ASSIGN_UNEXPECTED_REST       SyntaxErrorType = 152
	NO_LOCAL_VARIABLE                  SyntaxErrorType = 153
	NOT_EXPRESSION                     SyntaxErrorType = 154
	NUMBER_LITERAL_UNDERSCORE          SyntaxErrorType = 155
	NUMBERED_PARAMETER_IT              SyntaxErrorType = 156
	NUMBERED_PARAMETER_ORDINARY        SyntaxErrorType = 157
	NUMBERED_PARAMETER_OUTER_SCOPE     SyntaxErrorType = 158
	OPERATOR_MULTI_ASSIGN              SyntaxErrorType = 159
	OPERATOR_WRITE_ARGUMENTS           SyntaxErrorType = 160
	OPERATOR_WRITE_BLOCK               SyntaxErrorType = 161
	PARAMETER_ASSOC_SPLAT_MULTI        SyntaxErrorType = 162
	PARAMETER_BLOCK_MULTI              SyntaxErrorType = 163
	PARAMETER_CIRCULAR                 SyntaxErrorType = 164
	PARAMETER_METHOD_NAME              SyntaxErrorType = 165
	PARAMETER_NAME_REPEAT              SyntaxErrorType = 166
	PARAMETER_NO_DEFAULT               SyntaxErrorType = 167
	PARAMETER_NO_DEFAULT_KW            SyntaxErrorType = 168
	PARAMETER_NUMBERED_RESERVED        SyntaxErrorType = 169
	PARAMETER_ORDER                    SyntaxErrorType = 170
	PARAMETER_SPLAT_MULTI              SyntaxErrorType = 171
	PARAMETER_STAR                     SyntaxErrorType = 172
	PARAMETER_UNEXPECTED_FWD           SyntaxErrorType = 173
	PARAMETER_WILD_LOOSE_COMMA         SyntaxErrorType = 174
	PATTERN_EXPRESSION_AFTER_BRACKET   SyntaxErrorType = 175
	PATTERN_EXPRESSION_AFTER_COMMA     SyntaxErrorType = 176
	PATTERN_EXPRESSION_AFTER_HROCKET   SyntaxErrorType = 177
	PATTERN_EXPRESSION_AFTER_IN        SyntaxErrorType = 178
	PATTERN_EXPRESSION_AFTER_KEY       SyntaxErrorType = 179
	PATTERN_EXPRESSION_AFTER_PAREN     SyntaxErrorType = 180
	PATTERN_EXPRESSION_AFTER_PIN       SyntaxErrorType = 181
	PATTERN_EXPRESSION_AFTER_PIPE      SyntaxErrorType = 182
	PATTERN_EXPRESSION_AFTER_RANGE     SyntaxErrorType = 183
	PATTERN_EXPRESSION_AFTER_REST      SyntaxErrorType = 184
	PATTERN_HASH_KEY                   SyntaxErrorType = 185
	PATTERN_HASH_KEY_LABEL             SyntaxErrorType = 186
	PATTERN_IDENT_AFTER_HROCKET        SyntaxErrorType = 187
	PATTERN_LABEL_AFTER_COMMA          SyntaxErrorType = 188
	PATTERN_REST                       SyntaxErrorType = 189
	PATTERN_TERM_BRACE                 SyntaxErrorType = 190
	PATTERN_TERM_BRACKET               SyntaxErrorType = 191
	PATTERN_TERM_PAREN                 SyntaxErrorType = 192
	PIPEPIPEEQ_MULTI_ASSIGN            SyntaxErrorType = 193
	REGEXP_ENCODING_OPTION_MISMATCH    SyntaxErrorType = 194
	REGEXP_INCOMPAT_CHAR_ENCODING      SyntaxErrorType = 195
	REGEXP_INVALID_UNICODE_RANGE       SyntaxErrorType = 196
	REGEXP_NON_ESCAPED_MBC             SyntaxErrorType = 197
	REGEXP_TERM                        SyntaxErrorType = 198
	REGEXP_UTF8_CHAR_NON_UTF8_REGEXP   SyntaxErrorType = 199
	RESCUE_EXPRESSION                  SyntaxErrorType = 200
	RESCUE_MODIFIER_VALUE              SyntaxErrorType = 201
	RESCUE_TERM                        SyntaxErrorType = 202
	RESCUE_VARIABLE                    SyntaxErrorType = 203
	RETURN_INVALID                     SyntaxErrorType = 204
	SINGLETON_FOR_LITERALS             SyntaxErrorType = 205
	STATEMENT_ALIAS                    SyntaxErrorType = 206
	STATEMENT_POSTEXE_END              SyntaxErrorType = 207
	STATEMENT_PREEXE_BEGIN             SyntaxErrorType = 208
	STATEMENT_UNDEF                    SyntaxErrorType = 209
	STRING_CONCATENATION               SyntaxErrorType = 210
	STRING_INTERPOLATED_TERM           SyntaxErrorType = 211
	STRING_LITERAL_EOF                 SyntaxErrorType = 212
	STRING_LITERAL_TERM                SyntaxErrorType = 213
	SYMBOL_INVALID                     SyntaxErrorType = 214
	SYMBOL_TERM_DYNAMIC                SyntaxErrorType = 215
	SYMBOL_TERM_INTERPOLATED           SyntaxErrorType = 216
	TERNARY_COLON                      SyntaxErrorType = 217
	TERNARY_EXPRESSION_FALSE           SyntaxErrorType = 218
	TERNARY_EXPRESSION_TRUE            SyntaxErrorType = 219
	UNARY_RECEIVER                     SyntaxErrorType = 220
	UNDEF_ARGUMENT                     SyntaxErrorType = 221
	UNEXPECTED_TOKEN_CLOSE_CONTEXT     SyntaxErrorType = 222
	UNEXPECTED_TOKEN_IGNORE            SyntaxErrorType = 223
	UNTIL_TERM                         SyntaxErrorType = 224
	VOID_EXPRESSION                    SyntaxErrorType = 225
	WHILE_TERM                         SyntaxErrorType = 226
	WRITE_TARGET_IN_METHOD             SyntaxErrorType = 227
	WRITE_TARGET_READONLY              SyntaxErrorType = 228
	WRITE_TARGET_UNEXPECTED            SyntaxErrorType = 229
	XSTRING_TERM                       SyntaxErrorType = 230
)

var SyntaxErrorLevels = []SyntaxErrorLevel{SyntaxErrorFatal, SyntaxErrorArgument}
var SyntaxErrorTypes = []SyntaxErrorType{
	ALIAS_ARGUMENT,
	AMPAMPEQ_MULTI_ASSIGN,
	ARGUMENT_AFTER_BLOCK,
	ARGUMENT_AFTER_FORWARDING_ELLIPSES,
	ARGUMENT_BARE_HASH,
	ARGUMENT_BLOCK_FORWARDING,
	ARGUMENT_BLOCK_MULTI,
	ARGUMENT_FORMAL_CLASS,
	ARGUMENT_FORMAL_CONSTANT,
	ARGUMENT_FORMAL_GLOBAL,
	ARGUMENT_FORMAL_IVAR,
	ARGUMENT_FORWARDING_UNBOUND,
	ARGUMENT_IN,
	ARGUMENT_NO_FORWARDING_AMP,
	ARGUMENT_NO_FORWARDING_ELLIPSES,
	ARGUMENT_NO_FORWARDING_STAR,
	ARGUMENT_SPLAT_AFTER_ASSOC_SPLAT,
	ARGUMENT_SPLAT_AFTER_SPLAT,
	ARGUMENT_TERM_PAREN,
	ARGUMENT_UNEXPECTED_BLOCK,
	ARRAY_ELEMENT,
	ARRAY_EXPRESSION,
	ARRAY_EXPRESSION_AFTER_STAR,
	ARRAY_SEPARATOR,
	ARRAY_TERM,
	BEGIN_LONELY_ELSE,
	BEGIN_TERM,
	BEGIN_UPCASE_BRACE,
	BEGIN_UPCASE_TERM,
	BEGIN_UPCASE_TOPLEVEL,
	BLOCK_PARAM_LOCAL_VARIABLE,
	BLOCK_PARAM_PIPE_TERM,
	BLOCK_TERM_BRACE,
	BLOCK_TERM_END,
	CANNOT_PARSE_EXPRESSION,
	CANNOT_PARSE_STRING_PART,
	CASE_EXPRESSION_AFTER_CASE,
	CASE_EXPRESSION_AFTER_WHEN,
	CASE_MATCH_MISSING_PREDICATE,
	CASE_MISSING_CONDITIONS,
	CASE_TERM,
	CLASS_IN_METHOD,
	CLASS_NAME,
	CLASS_SUPERCLASS,
	CLASS_TERM,
	CLASS_UNEXPECTED_END,
	CONDITIONAL_ELSIF_PREDICATE,
	CONDITIONAL_IF_PREDICATE,
	CONDITIONAL_PREDICATE_TERM,
	CONDITIONAL_TERM,
	CONDITIONAL_TERM_ELSE,
	CONDITIONAL_UNLESS_PREDICATE,
	CONDITIONAL_UNTIL_PREDICATE,
	CONDITIONAL_WHILE_PREDICATE,
	CONSTANT_PATH_COLON_COLON_CONSTANT,
	DEF_ENDLESS,
	DEF_ENDLESS_SETTER,
	DEF_NAME,
	DEF_NAME_AFTER_RECEIVER,
	DEF_PARAMS_TERM,
	DEF_PARAMS_TERM_PAREN,
	DEF_RECEIVER,
	DEF_RECEIVER_TERM,
	DEF_TERM,
	DEFINED_EXPRESSION,
	EMBDOC_TERM,
	EMBEXPR_END,
	EMBVAR_INVALID,
	END_UPCASE_BRACE,
	END_UPCASE_TERM,
	ESCAPE_INVALID_CONTROL,
	ESCAPE_INVALID_CONTROL_REPEAT,
	ESCAPE_INVALID_HEXADECIMAL,
	ESCAPE_INVALID_META,
	ESCAPE_INVALID_META_REPEAT,
	ESCAPE_INVALID_UNICODE,
	ESCAPE_INVALID_UNICODE_CM_FLAGS,
	ESCAPE_INVALID_UNICODE_LITERAL,
	ESCAPE_INVALID_UNICODE_LONG,
	ESCAPE_INVALID_UNICODE_TERM,
	EXPECT_ARGUMENT,
	EXPECT_EOL_AFTER_STATEMENT,
	EXPECT_EXPRESSION_AFTER_AMPAMPEQ,
	EXPECT_EXPRESSION_AFTER_COMMA,
	EXPECT_EXPRESSION_AFTER_EQUAL,
	EXPECT_EXPRESSION_AFTER_LESS_LESS,
	EXPECT_EXPRESSION_AFTER_LPAREN,
	EXPECT_EXPRESSION_AFTER_OPERATOR,
	EXPECT_EXPRESSION_AFTER_PIPEPIPEEQ,
	EXPECT_EXPRESSION_AFTER_QUESTION,
	EXPECT_EXPRESSION_AFTER_SPLAT,
	EXPECT_EXPRESSION_AFTER_SPLAT_HASH,
	EXPECT_EXPRESSION_AFTER_STAR,
	EXPECT_IDENT_REQ_PARAMETER,
	EXPECT_LPAREN_REQ_PARAMETER,
	EXPECT_RBRACKET,
	EXPECT_RPAREN,
	EXPECT_RPAREN_AFTER_MULTI,
	EXPECT_RPAREN_REQ_PARAMETER,
	EXPECT_STRING_CONTENT,
	EXPECT_WHEN_DELIMITER,
	EXPRESSION_BARE_HASH,
	FLOAT_PARSE,
	FOR_COLLECTION,
	FOR_IN,
	FOR_INDEX,
	FOR_TERM,
	HASH_EXPRESSION_AFTER_LABEL,
	HASH_KEY,
	HASH_ROCKET,
	HASH_TERM,
	HASH_VALUE,
	HEREDOC_TERM,
	INCOMPLETE_QUESTION_MARK,
	INCOMPLETE_VARIABLE_CLASS,
	INCOMPLETE_VARIABLE_CLASS_3_3_0,
	INCOMPLETE_VARIABLE_INSTANCE,
	INCOMPLETE_VARIABLE_INSTANCE_3_3_0,
	INVALID_CHARACTER,
	INVALID_ENCODING_MAGIC_COMMENT,
	INVALID_FLOAT_EXPONENT,
	INVALID_MULTIBYTE_CHAR,
	INVALID_MULTIBYTE_CHARACTER,
	INVALID_MULTIBYTE_ESCAPE,
	INVALID_NUMBER_BINARY,
	INVALID_NUMBER_DECIMAL,
	INVALID_NUMBER_HEXADECIMAL,
	INVALID_NUMBER_OCTAL,
	INVALID_NUMBER_UNDERSCORE,
	INVALID_PERCENT,
	INVALID_PRINTABLE_CHARACTER,
	INVALID_VARIABLE_GLOBAL,
	INVALID_VARIABLE_GLOBAL_3_3_0,
	IT_NOT_ALLOWED_NUMBERED,
	IT_NOT_ALLOWED_ORDINARY,
	LAMBDA_OPEN,
	LAMBDA_TERM_BRACE,
	LAMBDA_TERM_END,
	LIST_I_LOWER_ELEMENT,
	LIST_I_LOWER_TERM,
	LIST_I_UPPER_ELEMENT,
	LIST_I_UPPER_TERM,
	LIST_W_LOWER_ELEMENT,
	LIST_W_LOWER_TERM,
	LIST_W_UPPER_ELEMENT,
	LIST_W_UPPER_TERM,
	MALLOC_FAILED,
	MIXED_ENCODING,
	MODULE_IN_METHOD,
	MODULE_NAME,
	MODULE_TERM,
	MULTI_ASSIGN_MULTI_SPLATS,
	MULTI_ASSIGN_UNEXPECTED_REST,
	NO_LOCAL_VARIABLE,
	NOT_EXPRESSION,
	NUMBER_LITERAL_UNDERSCORE,
	NUMBERED_PARAMETER_IT,
	NUMBERED_PARAMETER_ORDINARY,
	NUMBERED_PARAMETER_OUTER_SCOPE,
	OPERATOR_MULTI_ASSIGN,
	OPERATOR_WRITE_ARGUMENTS,
	OPERATOR_WRITE_BLOCK,
	PARAMETER_ASSOC_SPLAT_MULTI,
	PARAMETER_BLOCK_MULTI,
	PARAMETER_CIRCULAR,
	PARAMETER_METHOD_NAME,
	PARAMETER_NAME_REPEAT,
	PARAMETER_NO_DEFAULT,
	PARAMETER_NO_DEFAULT_KW,
	PARAMETER_NUMBERED_RESERVED,
	PARAMETER_ORDER,
	PARAMETER_SPLAT_MULTI,
	PARAMETER_STAR,
	PARAMETER_UNEXPECTED_FWD,
	PARAMETER_WILD_LOOSE_COMMA,
	PATTERN_EXPRESSION_AFTER_BRACKET,
	PATTERN_EXPRESSION_AFTER_COMMA,
	PATTERN_EXPRESSION_AFTER_HROCKET,
	PATTERN_EXPRESSION_AFTER_IN,
	PATTERN_EXPRESSION_AFTER_KEY,
	PATTERN_EXPRESSION_AFTER_PAREN,
	PATTERN_EXPRESSION_AFTER_PIN,
	PATTERN_EXPRESSION_AFTER_PIPE,
	PATTERN_EXPRESSION_AFTER_RANGE,
	PATTERN_EXPRESSION_AFTER_REST,
	PATTERN_HASH_KEY,
	PATTERN_HASH_KEY_LABEL,
	PATTERN_IDENT_AFTER_HROCKET,
	PATTERN_LABEL_AFTER_COMMA,
	PATTERN_REST,
	PATTERN_TERM_BRACE,
	PATTERN_TERM_BRACKET,
	PATTERN_TERM_PAREN,
	PIPEPIPEEQ_MULTI_ASSIGN,
	REGEXP_ENCODING_OPTION_MISMATCH,
	REGEXP_INCOMPAT_CHAR_ENCODING,
	REGEXP_INVALID_UNICODE_RANGE,
	REGEXP_NON_ESCAPED_MBC,
	REGEXP_TERM,
	REGEXP_UTF8_CHAR_NON_UTF8_REGEXP,
	RESCUE_EXPRESSION,
	RESCUE_MODIFIER_VALUE,
	RESCUE_TERM,
	RESCUE_VARIABLE,
	RETURN_INVALID,
	SINGLETON_FOR_LITERALS,
	STATEMENT_ALIAS,
	STATEMENT_POSTEXE_END,
	STATEMENT_PREEXE_BEGIN,
	STATEMENT_UNDEF,
	STRING_CONCATENATION,
	STRING_INTERPOLATED_TERM,
	STRING_LITERAL_EOF,
	STRING_LITERAL_TERM,
	SYMBOL_INVALID,
	SYMBOL_TERM_DYNAMIC,
	SYMBOL_TERM_INTERPOLATED,
	TERNARY_COLON,
	TERNARY_EXPRESSION_FALSE,
	TERNARY_EXPRESSION_TRUE,
	UNARY_RECEIVER,
	UNDEF_ARGUMENT,
	UNEXPECTED_TOKEN_CLOSE_CONTEXT,
	UNEXPECTED_TOKEN_IGNORE,
	UNTIL_TERM,
	VOID_EXPRESSION,
	WHILE_TERM,
	WRITE_TARGET_IN_METHOD,
	WRITE_TARGET_READONLY,
	WRITE_TARGET_UNEXPECTED,
	XSTRING_TERM,
}

type SyntaxError struct {
	Message  string
	Location *Location
	Level    SyntaxErrorLevel
	Type     SyntaxErrorType
}

func NewSyntaxError(
	message string,
	location *Location,
	level SyntaxErrorLevel,
	errType SyntaxErrorType,
) *SyntaxError {
	return &SyntaxError{
		Message:  message,
		Location: location,
		Level:    level,
		Type:     errType,
	}
}
