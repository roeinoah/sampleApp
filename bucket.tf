resource "aws_s3_bucket" "eks" {
  bucket = "cves"
  acl    = "public-read"

  tags {
    Name        = "cves"
    Environment = "sample"
  }
}