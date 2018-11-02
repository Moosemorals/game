# Game

Basically a dungeon crawler. Lets see how far I get.

## First night

I'm using the (termbox-go)[https://github.com/nsf/termbox-go] library to
manage the screen stuff. So far, so good. It's got a really simple
API, but I can use that to build up higher level constructs. (See: drawString)

At the moment I've got a couple of useful types, sprites and levels.

Sprites represent things that move (i.e., the player, so far), so they've 
got a position, and a move method (that's aware of walls and the edge 
of the screen)

Levels represent the envionment (things that don't move). It's mostly
a 1d array of bit flags masqurading as a 2d array.

Currently, it's driven by the termbox event loop. I wait for a keypress,
hand it to the sprite to deal with, then redraw the world.

Plan: Think about how to store the environment. Do I want a flag for
every type of thing? Doors can be open, closed, hidden, secret, found.
Should those all be flags, or do I want some kind of class heirachy?
(Except Go doesn't do that).

How about a canPass(sprite) bool function? Walls just return false,
doors check their state. Then levels becomes and array of passables,
and the wall and door types implment that interface. Feels clunky.

I'll have a think about it.