// mysql charset https://dev.mysql.com/doc/refman/5.7/en/charset-charsets.html
// postgres charset https://www.postgresql.org/docs/9.6/static/multibyte.html
variable "charset" {
  description = "The charset value. See MySQL's Supported Character Sets and Collations and Postgres' Character Set Support for more details and supported values. Postgres databases only support a value of UTF8 at creation time."
  type        = string
  default     = null
}

variable "collation" {
  description = "The collation value. See MySQL's Supported Character Sets and Collations and Postgres' Collation Support for more details and supported values. Postgres databases only support a value of en_US.UTF8 at creation time."
  type        = string
  default     = null
}
