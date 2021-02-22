## -- Syntax guide when Creating GoBot2Flags levels with .r2f files --

A guide to assist the creation of levels.

### -- Walls --

Horizontal walls are displayed with: `+--+--+--+--+` (this would be an example for a 4 unit wide wall).

Vertical walls however are displayed with: 
```
|
+
|
+
|
+
```
(this would be for a 3 unit tall wall).

### -- Floor --

There are 3 available colours: Red (symbolised by `R`), Blue (symbolised by `B`) and Yellow (symbolised by `Y`).

If, however you wished to place a flag then you would append the floor colour with an 'F' e.g. a flag in a blue square would be marked as 'BF'.

### -- The robot --

The robot is displayed with `>`, `<`, `^` or `v`.

### -- Other --

If you are in a space that is without a possible square that you could be placed on (where two `+` symbols intersect and there is no wall) then you put a `.` symbol as a placeholder: 

Example where two `+` intersect and there is no wall:
```
+--+--+
|R  R |
+  .  +
|R> RF|
+--+--+
```
Example where two `+` intersect and there is a wall:
```
+--+--+
|R  R |
+--+--+
|R> RF|
+--+--+
```
Also there needs to be an odd number of lines.