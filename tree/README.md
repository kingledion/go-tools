# tree

Package tree implements a simple tree that can be built from and stored to a
row based data format, such as a relational database or csv file.

A tree is here defined as a graph having three properties:
  - a single root node with no inbound edges
  - all non-root nodes have exactly one inbound edge (parent)
  - any node may have any number of outbound edges (children)
The nodes of the tree are assumed to have a primary identifier or key by
which parent and child relationships can be defined.

## Objectives

The goal fo this tree is to implement a generic tree with minimal memory usage and maximal performance for general use cases. This package does not philosophically attempt to avoid panic-ing, as that safety comes at the price of sacrificed performance. Therefore, it would be expected that an implementer will perform nil checks as appropriate. 

## Development

Documentation can be greated with `godoc`.