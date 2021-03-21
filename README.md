# Advent Of Code - 2019
Advent of code 2nd attempt (2019) - https://adventofcode.com/

---
## Motivation 🚂
Following my experience with [Rust](Rust) a couple of months before I wanted to keep on experimenting with another programming languages.
Creativity was never my strong suit, so I never know what do to in order to experiment with new languages, by solving Advent of Code I do not have to worry with the step of creating a challenge and go straight into solving them.

I decided to start this edition of AOC so close to the last one mainly because of two points:
- I had already a language in mind to use to solve this edition (Spoilers: [Go](Go))
- I had almost a month of free time between semester, and dduring the semester is definetely not the time to feel "pressured" into solving a puzzle a day (and in the days that I am not able having it stuck in my mind until I do)

---

## Why Go 🐭 ( No Gopher emoji I think )
To be fair I had already done a really small project in [Go](Go) a couple of years ago, it was not a full project, not even close to that, but I had done a small change to a colleague's project in [Go](Go), in order to adapt it for a slightly different purpose that I had.

From the previous interaction that I had with [Go](Go) I knew that it was not a particularry difficult language to get into, to do something quick and easy, but I had no deeper understanding of the overall language, that is why I chose it!

---

## Methodology ✍
I have to admit that having just come from experimenting with [Rust](Rust) by comparisson [Go](Go) felt a lot less ... pleasing. Yes just like I mentioned [Go](Go) seems like a simple language to do something quick for a specific purpose, but when it comes to making sure that the code runs as expected, the types and the overall feel, in my (personal) opinion felt less pleasing.

The fact that the formatting as to be a certain way annoys me a bit, for example I wanted:
```
if (condition) { action }

if (condition) action
```

But I was forced to do:
```
if (condition) {
    action
}
```
This got especially annoying around simple if - else block that could be easily replaced by ternary operator (which does not exist in [Go](Go)).

The documentation around [Go](Go) is quick and easy so there was never any big difficulty just searching by the specific problem I was having and understanding the solution (so that is definitely a pro when comparing to [Rust](Rust)).

For the methodology of solving the problems I find important to point out that I was not able to solve all of them in order. Of course I tried to, but there were two different days in which I could not for the life of me solve the second part, not even find a plausible methodology to solve the problem so I decided to skip just like I said I would try to do in the previous edition.

---

## Most interesting challenges 💪

- ### [Day 02](https://adventofcode.com/2019/day/2)

Yup straight away on Day 02, altough this is solely my fault, for some reason I thought that for Part 2 checking all the possible combinations would take too long ... I was wrong. It definitely was a teaching moment ... In deed sometimes the easiest solution is the right one.

- ### [Day 05](https://adventofcode.com/2019/day/5)

The first of many interesting IntCode Program challenges. Nothing to specific to talk about in this one, but I found this group of challenges really interesting, it reminded me of a subject I had in my course.

- ### [Day 08](https://adventofcode.com/2019/day/8)

Once again, not particulary hard but I loved solving the problem exporting the image and being able to read the solution in the image. For some reason having the solution by its own in a file made it more rewarding?

- ### [Day 17](https://adventofcode.com/2019/day/17)

I absolutely loved this challenge, trying to express to a program a way to divide a set of commands in N sub commands which could then be repeated because memory was limited. What made this challenge so interesting to finish was mainly do two reasons:
1. Not even I knew how I (a human) would solve the challenge in a programatically way, so finding a solution for a computer which needed explicit rules felt rewarding and meaningful.
2. Apparently there was an algorithm which is able to solve this kind of challenges, I did not know that algorithm, neither did I need it.

- ### [Day 18](https://adventofcode.com/2019/day/18)

For some reason this challenge took me so much time, I was stuck on it for so long. The first problem was that I was not dealing with it as a graph problem once I understood that my code was already a mess from so many different attempts. So when I finally completed Part 1 I was not in the right mental space to solve Part 2.

When I finally came around to this challenge once again the first thing I did was fully delete what I had done and start over (the right decision without a doubt). After doing this solving Part 2 was easy.

- ### [Day 21](https://adventofcode.com/2019/day/21)

This one was intringuing, to be honest I have no idea how would one solve this challenge programatically. What I did was code it in a way that would take one file as input and then change this file (mannualy and iteratively) adapting to the situations in which the code fails, eventually coming across a valid possibility.

- ### [Day 22](https://adventofcode.com/2019/day/22)

For Part 2 of this challenge I was also stuck, I could not for the life of me solve it in a way that it would not simply take to long. Eventually I had to search for a hint.

I noticed that some people had solved it by finding a linear function which described the suffle operations and consequently one that would decribe the whole shuffle, then this function could be aggregated to itself the number of times needed. This was enough for me to understand how to solve and eventually got there, by the end I felt quite satisfied with the code that I had done and with the puzzle itself.

---

## What comes next 🔮
To be honest I have another projects in mind for the near future as well as having the current semester to survive to. Because of that I am not sure if the next time I take on a AOC challenge will be the edition of 2018 or the edition of 2021, if it will be with yet another new language or one already used in another edition.

Just like with erveything else ... only time will tell.

---

## Acknowledgments 🤗

Once again the challenges developed by [Eric Wastl][Eric Wastl Webpage] were a lot of fun to solve. This year the IntCode Program challenges were so interesting, the way you had to build on code that you already had done before, amazed me. Of course the input could simply be given to us, but being given with this extra added step made it more pleasing. Surely this could easily become to much, I do not wish for the 25 puzzles to become a single and longer one, but it seemed like a right mix between IntCode and not IntCode challenges.

Altough I have to admit I found this edition as whole a bit harder than the edition of 2020, but nothing wrong with that.

[Go]: https://golang.org/
[Rust]: https://www.rust-lang.org/
[Eric Wastl Webpage]: http://was.tl/