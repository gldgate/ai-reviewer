---
id: proto-generated
model_category: cheap
path_filters:
  - "**/*.pb.go"
  - "**/*.pulsar.go"
  - "**/*.proto"
---
These files are auto-generated from Protocol Buffer (.proto) definitions using the `protoc` compiler and should not be manually edited. Any changes to these files will be overwritten the next time the code is generated.

You will see changes to them in PRs, but only as updates. Anything that looks like custom code should be called out in a review.

If there are changes in .proto files, there should be corresponding changes in the generated code. If they are not, or they clearly do not match, warn the user to run `ignite generate proto-go` again.