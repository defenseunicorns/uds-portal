
###
# DO NOT DELETE OR TEARDOWN
# WILL BE USED FOR NIGHTLY DEMO DNS RECORDS
# MUST REMAIN STATIC
###
resource "aws_eip" "nightly_eip" {
  tags = {
    Name        = "runtime-nightly"
  }
}
