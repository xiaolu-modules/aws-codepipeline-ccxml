variable "bucket" {
  description = "The bucket that will contain the feed"
  type        = string
}

variable "key" {
  description = "The key within the bucket that will contain the feed"
  type        = string
  default     = "cc.xml"
}

variable "function_name" {
  description = "The name of the Lambda function"
  type        = string
  default     = "ccxml"
}

variable "memory_size" {
  description = "The amount of memory to allocate to the Lambda function"
  type        = number
  default     = 128
}

variable "timeout" {
  description = "The timeout for the Lambda function in seconds"
  type        = number
  default     = 20
}

variable "tags" {
  description = "A map of tags to add to all resources"
  type        = map(string)
  default     = {}
}
