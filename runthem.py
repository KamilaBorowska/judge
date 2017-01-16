from itertools import product
from collections import defaultdict
from subprocess import call

programs = [
    'python prog1.py',
    'python -c print("PONG\\n7")',
]

sizes = list(range(3, 10))

winners = {}

for program in programs:
    winners[program] = 0

for params in product(sizes, programs, programs):
    [size, prog1, prog2] = params
    if prog1 == prog2:
        continue
    print('\njudge {} "{}" "{}"'.format(size, prog1, prog2))
    winner = call(['judge', '--no-graphics', str(size), prog1, prog2])
    if winner == 0:
        raise "oh no"
    winners[params[winner]] += 1

print("\n\n\nWYGRALI:")
battles = len(sizes) * (len(programs) - 1) * 2
for [program, wins] in reversed(sorted(winners.items(), key=lambda v: v[1])):
    print("{:4} {:>7.2%} {}".format(wins, float(wins) / battles, program))
print("Każdy gracz miał {} potyczek".format(battles))