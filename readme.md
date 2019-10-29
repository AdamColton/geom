## Geom
A vector package for Go. The two main sub packages are d2 and d3 for 2D and 3D
math.

### Parametric Naming Convention
Parametric methods take floating point arguments labeled as t0, t1, ..., tn. 

The naming convention for parametric methods is as follows. First is the return
type, either Pt or V. Then the number of parametric arguments. If it is a
currying method, this will be followed by c and the number of arguments curried
by the call.

So Pt1 takes one parametric argument and returns a Pt. Pt2c1 curries one of two
parametric aruments, returning a Pt1 interface.

The arguments are anchored to the range of [0, 1]. So for a curve t0=0
represents the starting point and t0=1 represents the ending point. In the case
of a curve, it may be meaningful to use values outside this range, for instance,
on a line t0=-1 is a point equal distance from the end in the opposite
direction. In the case of a surface, values outside this range may not be
meaninful.