# Software Requirements Specification for Spring 2020 "Software Engineering" Primitives Project at Wright State University

## Overview

This document specifies the project modifation to Fogleman's Primitive Application, hereafter called "App".

## References

## Functional Requirements

### [User Story 2A ](features.md "Ref. Features And User Stories")


Req. 2.0 The App shall for each output image identify each distinct color within the output image using a different number per color.

Req. 2.1 The App shall annotate each output image to show the number corresponding to each area of distict color in the output.

Req. 2.1.1 Each annotation shall consist of text in 8pt font or greater.

Req. 2.1.2 Each annotation shall be presented in each output image directly above the distict color to which the annotation corresponds.

Req. 2.1.2.1 Each area of distinct color in the output shall have a corresponding annotation.

Req. 2.2 The App shall output an image consisting of outlines for each area of distict color.
 
Req. 2.2.1 Output outlines shall be annotated in the same manner as specified for areas of distict color in Req. 2.1.

Req. 2.3 The App shall output a table correlating numbers with distinct colors that appears in the output image.

Req. 2.4 The App shall have levels of difficulty to limit the number of geometric shapes annotated.

Req. 2.4.1 The App shall have an easy level consisting of 50 geometric shapes.

Req. 2.4.2 The App shall have a medium level consisting of 100 geometric shapes.

Req. 2.4.3 The App shall have a hard level consisting of 150 geometric shapes. 

Req. 2.4.4 The App shall have an expert level consisting of 300 geometric shapes. 

### [User Story 3 ](features.md "Ref. Features And User Stories")

Req. 3.0 The User shall have an option to apply a filter to the output image. 

Req. 3.1 The App shall produce an output image in gray scale.

Req. 3.2 The App shall produce an output image in sepia. 

Req. 3.3 The App shall produce an output image in negative.

Req. 3.4 The User shall be able to select a filter.

### [User Story 4 ](features.md "Ref. Features And User Stories")

Req. 4.0 The User Interface shall have a help button to walk the user through how to use primitive.

Req. 4.1 The User Interface shall have an export option to export the image as a paint by numbers.

Req. 4.2 The User Interface shall allow the user to select the number of geometric shapes to form their image.

Req. 4.3 The User Interface shall allow the user to select up to five different geometric shapes.

Req. 4.4 The User Interface shall display the final image.

Req. 4.5 The User shall be able to select an image from their device to modify.
