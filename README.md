# Gorup

Randomly picks a name from a list. Use `-q` to control the chance to pick each name.

## Usage

```
gorup apple pear banana
```

Pick a random name among 'apple', 'pear', and 'banana', each having an equal chance of being picked (=1:1:1)

```
gorup cat -q2 dog -q 0.3 mouse pig
```

Pick a random name among 'cat', 'dog', 'mouse', and 'pig'. Chance to pick is _cat_ : _dog_ : _mouse_ : _pig_ = 1 : 2 : 0.3 : 1.
