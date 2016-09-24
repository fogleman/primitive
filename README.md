# Primitive Pictures

Reproducing images with geometric primitives.

![Example](https://www.michaelfogleman.com/static/primitive/examples/16550611738.200.128.4.5.png)

### How it Works

A target image is provided as input. The algorithm tries to find the single most optimal shape that can be drawn to minimize the error between the target image and the drawn image. It repeats this process, adding *one shape at a time*. Around 50 to 200 shapes are needed to reach a result that is recognizable yet artistic and abstract.

### How it Works, Part II

Say we have a `Target Image`. This is what we're working towards recreating. We start with a blank canvas, but we fill it with a single solid color. Currently, this is the average color of the `Target Image`. We call this new blank canvas the `Current Image`. Now, we start evaluating shapes. To evaluate a shape, we draw it on top of the `Current Image`, producing a `New Image`. This `New Image` is compared to the `Target Image` to compute a score. We use the [root-mean-square error](https://en.wikipedia.org/wiki/Root-mean-square_deviation) for the score.

    Current Image + Shape => New Image
    RMSE(New Image, Target Image) => Score

The shapes are generated randomly. We can generate a random shape and score it. Then we can mutate the shape (by tweaking a triangle vertex, tweaking an ellipse radius or center, etc.) and score it again. If the mutation improved the score, we keep it. Otherwise we rollback to the previous state. Repeating this process is known as [hill climbing](https://en.wikipedia.org/wiki/Hill_climbing). Hill climbing is prone to getting stuck in local minima, so we actually do this many different times with several different starting shapes. We can also generate N random shapes and pick the best one before we start hill climbing. [Simulated annealing](https://en.wikipedia.org/wiki/Simulated_annealing) is another good option, but in my tests I found the hill climbing technique just as good and faster, at least for this particular problem.

Once we have found a good-scoring shape, we add it to the `Current Image`, where it will remain unchanged. Then we start the process again to find the next shape to draw. This process is repeated as many times as desired.

### Twitter

Follow [@PrimitivePic](https://twitter.com/PrimitivePic) on Twitter to see a new primitive picture every 30 minutes!

The Twitter bot looks for interesting photos using the Flickr API, runs the algorithm using randomized parameters, and
posts the picture using the Twitter API.

You can tweet a picture to the bot and it will process it for you.

### Command-line Usage

Run it on your own images! First, [install Go](https://golang.org/doc/install).

    go get -u github.com/fogleman/primitive
    primitive -i input.png -o output.png -n 100

Small input images should be used (like 256x256px). You don't need the detail anyway and the code will run faster.

| Flag | Default | Description |
| --- | --- | --- |
| `-i` | n/a | input file |
| `-o` | n/a | output file |
| `-n` | n/a | number of shapes |
| `-m` | 1 | mode: 0=combo, 1=triangle, 2=rect, 3=ellipse, 4=circle, 5=rotatedrect |
| `-r` | 256 | resize large input images to this size before processing |
| `-s` | 1024 | output image size |
| `-a` | 128 | color alpha |
| `-v` | off | verbose output |

### Output Formats

Depending on the output filename extension provided, you can produce different types of output.

- `PNG`: raster output
- `JPG`: raster output
- `SVG`: vector output
- `GIF`: animated output showing shapes being added - requires ImageMagick (specifically the `convert` command)

For PNG and SVG outputs, you can also include `%d`, `%03d`, etc. in the filename. In this case, each frame will be saved separately.

You can use the `-o` flag multiple times. This way you can save both a PNG and an SVG, for example.

### Primitives

The following primitives are supported:

- Triangle
- Rectangle (axis-aligned)
- Ellipse (axis-aligned)
- Circle
- Rotated Rectangle
- Combo (a mix of the above in a single image)

More shapes can be added by implementing the following interface:

```go
type Shape interface {
	Rasterize() []Scanline
	Copy() Shape
	Mutate()
	Draw(dc *gg.Context)
	SVG(attrs string) string
}
```

### Features

- [Hill Climbing](https://en.wikipedia.org/wiki/Hill_climbing) or [Simulated Annealing](https://en.wikipedia.org/wiki/Simulated_annealing) for optimization (hill climbing multiple random shapes is nearly as good as annealing and faster)
- Scanline rasterization of shapes in pure Go (preferable for implementing the features below)
- Optimal color computation based on affected pixels for each shape (color is directly computed, not optimized for)
- Partial image difference for faster scoring (only pixels that change need be considered)
- Anti-aliased output rendering

### Inspiration

This project was originally inspired by the popular and excellent work of Roger Johansson - [Genetic Programming: Evolution of Mona Lisa](https://rogeralsing.com/2008/12/07/genetic-programming-evolution-of-mona-lisa/). Since seeing that article when it was quite new, I've tinkered with this problem here and there over the years. But only now am I satisfied with my results.

It should be noted that there are significant differences in my implementation compared to Roger's original work. Mine is not a genetic algorithm. Mine only operates on one shape at a time. Mine is much faster (AFAIK) and supports many types of shapes.

### Progression

This GIF demonstrates the iterative nature of the algorithm, attempting to minimize the mean squared error by adding one shape at a time. (Use a ".gif" output file to generate one yourself!)

![Mona Lisa](https://www.michaelfogleman.com/static/primitive/examples/monalisa.gif)

### Static Animation

Since the algorithm has a random component to it, you can run it against the same input image multiple times to bring life to a static image.

![Pencils](https://www.michaelfogleman.com/static/primitive/examples/pencils.gif)

### Creative Constraints

If you're willing to dabble in the code, you can enforce constraints on the shapes to produce even more interesting results. Here, the rectangles are constrained to point toward the sun in this picture of a pyramid sunset.

![Pyramids](https://www.michaelfogleman.com/static/primitive/examples/pyramids.png)

### Shape and Iteration Comparison Matrix

The matrix below shows triangles, ellipses and rectangles at 50, 100 and 200 iterations each.

![Matrix](http://i.imgur.com/H5NYpL4.png)

### Examples

Here are more examples from interesting photos found on Flickr.

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
![Example](https://www.michaelfogleman.com/static/primitive/examples/27471731151.50.128.4.1.png)
![Example](https://www.michaelfogleman.com/static/primitive/examples/11720700033.200.128.4.3.png)
![Example](https://www.michaelfogleman.com/static/primitive/examples/18782606664.png)
![Example](https://www.michaelfogleman.com/static/primitive/examples/21374478713.png)
![Example](https://www.michaelfogleman.com/static/primitive/examples/15196426112.200.128.4.5.png)
![Example](https://www.michaelfogleman.com/static/primitive/examples/24696847962.png)
![Example](https://www.michaelfogleman.com/static/primitive/examples/18276676312.100.128.4.1.png)
