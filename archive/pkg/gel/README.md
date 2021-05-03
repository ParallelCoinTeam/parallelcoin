# gel

Gio elements

This is a widget toolkit based on the basic set of material widgets built with Gio.

So far, the following additional widgets have been created - scrollbar for the scrolling list, a column and basic table,
and a top-bar side-bar status-bar multi-page application framework, and a background process runner to minimise
interruptions to the rendering of the interfaces.

There is also a state-widget pool generated with the stuff in `poolgen/` that you initialize before starting the render
loop that can then be used anywhere to add new clickables, bools, inputs and so on, without having to double up with a 
pre-specification, they allocate first run and then cache thereafter, maintaining their state.

Gel uses fluent programming techniques to simplify and denoise the visual structure of Gio widget definitions. To take
maximum advantage of it, make use of the fact that Go allows breaking lines after dot operators between the chained
methods such that each movable section can be pulled up and moved into another position without dealing with a screenful
of red underlines in your IDE caused by unpaired brackets.

Future plans include creating a serialization format for events and render op queues and accompanying socket transports
to enable remote display and control of user interfaces.
