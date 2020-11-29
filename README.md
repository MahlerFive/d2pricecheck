d2pricecheck
============

A utility to analyze posted prices for items in Diablo 2.

Prerequisites
-------------

1. Install [Go](https://golang.org/).
2. Install git (optional if you want to use `git clone` to get the repo).

How to Use
----------

1. Download / clone the repo. If using `git`:

```
$ git clone git@github.com:MahlerFive/d2pricecheck.git
```

2. Prepare an input file.

Copy a bunch of text from a channel where people are posting items for trade. For example, in the [Project Diablo 2 discord](https://discord.gg/RgX4MWu) channel `#sc-trade-💰`.

Now, paste the text into `input.txt` in the `d2pricecheck` folder.

3. Run the price analyzer.

```
$ cd [path to d2pricecheck folder]
$ go run cmd/priceanalyzer/main.go --in input.txt --out output.txt
```

If you don't provide the `--in` or `--out` flags, it will use `input.txt` and `output.txt` by default. Make sure the path to the files is relative to the `d2pricecheck` folder.

4. Search the output file.

Open `output.txt` or whatever file you specified as the `--out` parameter. Search the file for the item you want to price check.

For example, if you search for `Guardian Angel` you might see something like:

```
Guardian Angel  um:1    mal:2
```

This means there was one person pricing Guardian Angel at um rune, and two people pricing it as a mal rune.

Limitations
-----------

There are a lot of limitations on the text parsing right now.

- Only works for unique and set items
- Item name and price must be on the same line
- Name must exactly match the full name
- Item must be first, and the price must be at the end
- Only rune prices considered
- Doesn't work when there are other letters other than the item name and the rune name (eg. specifying modifiers like `eth`/`ethereal` or `180ED`, starting the line with `[o]` to indicate an offer)

<h4>Examples that work</h4>

Tal Rasha's Lidless Eye, IST

Tal Rasha's Lidless Eye=IST

Lidless Wall [+1] =pul

<h4>Examples that won't work</h4>

Tal weapon, IST

Tal Rasha's Lidless Eye, IST/GUL

Eth Lidless Wall [+1] =pul

[O] Lidless Wall [+1] =pul

Lidless Wall +2 Skills, IST

Future
------

This is just a hobby project so I may or may not get to these.

- Improve text parsing to remove most of the existing limitations
- Sorted rune price output
- Output examples of each price (this should help you see how different modifiers affect the price)
- Support more item types (eg. Socket item bases, common magic items like 15/40 jewels or 3/20 javs)
- Automate input by scraping Discord with a bot
- Provide a web interface and/or Discord bot for searching

Note that the intention of this project is just to help with easily checking prices so you know what to prices to post trades for, what to offer for trades, and have an idea of what items are worth keeping/muling. Even if a site is built, it's not intended to be a trading platform. This is both to limit scope and because I'm building this with the intention of supporting the [Project Diablo 2](https://www.projectdiablo2.com/) community, where external trading sites are prohibited.

Contributing
------------

**See any text parsing limitations that aren't listed above?**

Please post a Github issue with the following info:

- Description of what text doesn't get parsed properly
- At least one example
- Tag with `bug`

**Want to make your own trade prices searchable?**

Please use an easy to parse format.
- Use full item names
- Put rune price at the end of the line separated by a space or punctuation
- Avoid using any other letters in the description other than the item name or rune name

See the examples in the limitations section.

**Feature Requests**

Please post a Github issue and tag it with `enhancement`.

**Code Contributions**

At the moment, I'm not accepting code contributions since this is a very early version and the code is in major flux. As things settle I will open it up.
