input_file = 'words.csv'

go_lines = [
    '// !!!           WARNING           !!!',
    '// !!! THIS FILE IS AUTO-GENERATED !!!',
    '// !!!         DO NOT EDIT         !!!',
    '',
    'package translation',
    '',
    'var words = []Word{'
]

file = open(input_file, encoding='utf-8')
lines = file.readlines()

counter = 0
for line in lines:
    counter += 1
    parts = line.split(';')
    if len(parts) < 3:
        continue

    spanish = parts[0].strip()
    pos = parts[1].strip()
    english = parts[2].strip()
    go_lines.append(f'\t{{ ID: {counter}, Spanish: "{spanish}", English: "{english}", PoS: "{pos}" }},')

go_lines.append('}')

go_code = '\n'.join(go_lines)
print(go_code)
