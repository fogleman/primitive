# Primitive Pictures

Reproducing images with geometric primitives.

![Example](https://www.michaelfogleman.com/static/primitive/examples/16550611738.200.128.4.5.png)

### How it Works

A target image is provided as input. The algorithm tries to find a shape that can be drawn to minimize the error
between the target image and the drawn image. It repeats this process, adding one shape at a time.

### Twitter

Follow [@PrimitivePic](https://twitter.com/PrimitivePic) on Twitter to see a new primitive picture every 30 minutes!

The Twitter bot looks for interesting photos using the Flickr API, runs the algorithm using randomized parameters, and
posts the picture using the Twitter API.

### Command-line Usage

    go get -u github.com/fogleman/primitive
    primitive -i input.png -o output.png -n 100

| Flag | Default | Description |
| --- | --- | --- |
| -i | n/a | input file |
| -o | n/a | output file |
| -n | n/a | number of shapes |
| -m | 1 | mode: 0=combo, 1=triangle, 2=rect, 3=ellipse, 4=circle, 5=rotatedrect |
| -s | 1 | output scaling factor |
| -a | 128 | color alpha |

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

### Inspiration

This project was originally inspired by the popular and excellent work of Roger Johansson - [Genetic Programming: Evolution of Mona Lisa](https://rogeralsing.com/2008/12/07/genetic-programming-evolution-of-mona-lisa/). Since seeing that article when it was quite new, I've tinkered with this problem here and there over the years. But only now am I satisfied with my results.

### Progression

This GIF demonstrates the iterative nature of the algorithm, attempting to minimize the mean squared error by adding one shape at a time.

![Mona Lisa](https://www.michaelfogleman.com/static/primitive/examples/monalisa.gif)

### Static Animation

Since the algorithm has a random component to it, you can run it against the same input image multiple times to bring life to a static image.

![Pencils](https://www.michaelfogleman.com/static/primitive/examples/pencils.gif)

### Shape and Iteration Comparison Matrix

The matrix below shows triangles, ellipses and rectangles at 50, 100 and 200 iterations each. For a similar comparison of 500 different images, see: [large comparison matrix](https://www.michaelfogleman.com/static/primitive/) (warning: this page has 6000 images on it!)

![Matrix](http://i.imgur.com/H5NYpL4.png)

### Examples

Here are more examples from interesting photos found on Flickr.

![Example](https://www.michaelfogleman.com/static/primitive/examples/27471731151.50.128.4.1.png)
![Example](https://www.michaelfogleman.com/static/primitive/examples/11720700033.200.128.4.3.png)
![Example](https://www.michaelfogleman.com/static/primitive/examples/18782606664.png)
![Example](https://www.michaelfogleman.com/static/primitive/examples/21374478713.png)
![Example](https://www.michaelfogleman.com/static/primitive/examples/15196426112.200.128.4.5.png)
![Example](https://www.michaelfogleman.com/static/primitive/examples/24696847962.png)
![Example](https://www.michaelfogleman.com/static/primitive/examples/18276676312.100.128.4.1.png)
![Example](https://www.michaelfogleman.com/static/primitive/examples/29167683201.png)
![Example](https://www.michaelfogleman.com/static/primitive/examples/26574286221.200.128.4.1.png)
![Example](https://www.michaelfogleman.com/static/primitive/examples/15011768709.200.128.4.1.png)
![Example](https://www.michaelfogleman.com/static/primitive/examples/27540729075.200.128.4.1.png)
![Example](https://www.michaelfogleman.com/static/primitive/examples/28896874003.png)
![Example](https://www.michaelfogleman.com/static/primitive/examples/20414282102.png)
![Example](https://www.michaelfogleman.com/static/primitive/examples/15199237095.200.128.4.1.png)
![Example](https://www.michaelfogleman.com/static/primitive/examples/11707819764.200.128.4.1.png)
![Example](https://www.michaelfogleman.com/static/primitive/examples/18270231645.200.128.4.3.png)
![Example](https://www.michaelfogleman.com/static/primitive/examples/15705764893.png)
![Example](https://www.michaelfogleman.com/static/primitive/examples/25213252889.png)
![Example](https://www.michaelfogleman.com/static/primitive/examples/15015411870.200.128.4.3.png)
![Example](https://www.michaelfogleman.com/static/primitive/examples/25766500104.png)
