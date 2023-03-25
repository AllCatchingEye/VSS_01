//go:generate protoc -I=. --go_out=. --go_opt=paths=source_relative --proto_path=. messages/messages.proto
package generate
