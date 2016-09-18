# Primitive Pictures

Reproducing images with geometric primitives.

![Examples](http://i.imgur.com/H5NYpL4.png)

### Twitter

Follow [@PrimitivePic](https://twitter.com/PrimitivePic) on Twitter to see a new primitive picture every 30 minutes!

The Twitter bot looks for interesting photos using the Flickr API, runs the algorithm using randomized parameters, and
posts the picture using the Twitter API.

### How it Works

A target image is provided as input. The algorithm tries to find a shape that can be drawn to minimize the error
between the target image and the drawn image. It repeats this process, adding one shape at a time.

### Primitives

The following primitives are supported:

- Triangle
- Rectangle (axis-aligned)
- Ellipse (axis-aligned)
- Circle
- Rotated Rectangle
- Combo (a mix of the above in a single image)

### Features

- Hill Climbing or Simulated Annealing for optimization
- Optimal color computation based on affected pixels for each shape
- Partial image difference for a faster scoring function
- Anti-aliased output rendering
