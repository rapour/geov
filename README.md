# Visvalingam–Whyatt simplification for Multipolygons

Geov is a simple Go command-line tool to simplify a set of juxtaposed polygons known as trapezoidal maps without generating inconsistent borders between them. It can also be used as a library in a Go project.

![Simplified multipolygons](https://i.ibb.co/Lz7C0vD/simplified.png)

## Problem 

The [Visvalingam–Whyatt algorithm](https://bost.ocks.org/mike/simplify/) is a line simplification algorithm intended to simplify polylines. It performs well on a single polygon. The issue arises when one wants to simplify a bunch of polygons that share borders. If the algorithm is applied to each polygon individually, the ensuing simplified polygons will have mismatched borders as the algorithm acts differently on the shared segment of two juxtaposed polygon. The subsequent multipolygon is no longer a trapezoidal maps and its polygons may have overlapped with each other. 

## Solution 

As Mike Bostock has explained in this [blog post](https://bost.ocks.org/mike/topology/), an intuitive solution is to transform the multipolygon geometry to a topology. As opposed to geometry, a topology partitions the multipolygon to a set of unique polylines known as arcs. Therefore, the topology can be inferred in a way that the shared segment of two polygons becomes a single arc. The arc belongs to both the polylgons but now since there is only one polyline (arc) to simplify, the two polygons inherit the same simplified polyline and therefore maintain their borders with each other. 

## How Geov works as a library

Geov uses [S2](https://github.com/kellydunn/golang-geo) as the base library to work with polygons. Although this library specialized for geospatial points, it can be used to represent any geometry in the two dimensional space. 

**_IMPORTANT:_** The library assumes the first and last points in each polygon are the same. Maintain this convention or you might get inconsistent results. 

First you need to employ S2 to defined your multipolygon as a `[]*geo.Polygon` slice. Then simplifying these polygons is as easy as: 
```go

var polygonSlice []*geo.Polygon

/*
    populate polygonSlice with your multipolygon data...
*/

mp := NewMultiPolygon(polygonSlice)

simplifiedMultiPolygon := geov.Simplify(mp, geov.Visvalingam, ratio)

```

The `ratio` is a factor in `[0 1]` range that describes how simple you want your multipolygon to be. The ratio of one doesn't perform any simplification on the multipolygon and the ration of zero results in the simplest form of the polygon that maintains the topology. 

Alternatively, you could represent your multipolygon as a geojson of `FeatureCollection` type with arbitrary number of `Polygon` geometries and use the library's unmrashal method:

```go
bin, err := os.Create("path/to/geojson/file")
if err != nil {
    log.Fatal(err)
}

mp, err := geov.Unmarshal(bin)
if err != nil {
    log.Fatal(err)
}

simplifiedMultiPolygon := geov.Simplify(mp, geov.Visvalingam, ratio)

```

## Visualizing your Multipolygon

Geov can generate an SVG output of your multipolygons. You just need to pass an `io.Writer` to the visualization method to write the contents of the svg into it as follows:
```go
out, err := os.Create("path/to/output/file.svg")
		if err != nil {
			log.Fatal(err)
		}

err = simplifiedMultiPolygon.SVG(out)
		if err != nil {
			log.Fatal(err)
		}
```

## How Geov works as a command-line tool 

TODO

## Constribution

If you want to contribute to this project and make it better, your help is appreciated. Please open a pull request and let's discuss it further. 
