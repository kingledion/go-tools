# tree

Package tree implements a simple tree that can be built from and stored to a
row based data format, such as a relational database or csv file.

A tree is here defined as a graph having three properties:
  - a single root node with no inbound edges
  - all non-root nodes have exactly one inbound edge (parent)
  - any node may have any number of outbound edges (children)
The nodes of the tree are assumed to have a primary identifier or key by
which parent and child relationships can be defined.

## Development

Documentation can be greated with `godoc`