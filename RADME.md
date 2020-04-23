# Bencode

This module provides functionality to encode and decode bencode data.

Work in progres.

#### TODO
- decode:
    - [ ] into interface{} argument as general case
    - [x] into agrument of supported type (int, string, []T, map[string]T etc.)
    - [ ] into struct:
        - [ ] match field by it's name
        - [x] match field by tag
    - [x] as raw bencode
- encode:
    - [ ] argument of type T
    - [ ] struct
