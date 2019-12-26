cat bench.json > benchbig.json
for i in $(seq 20); do
  cp benchbig.json benchbignew.json
  cat benchbig.json >> benchbignew.json
  mv benchbignew.json benchbig.json
done
for i in $(seq 10); do
  cat bench.json >> benchmedium.json
done