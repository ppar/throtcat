
# throtcat

`throtcat` is a `cat(1)` -style program that throttles it throughput based 
on external criteria.

Two criteria for throttling are available:

- `time`: pause N milliseconds between each input line
- `innodb`: connect to a MySQL server, observe its InnoDB engine metrics and 
   pause when certain thresholds are reached


It can be used in a scenario such as 

```
$ throtcat [OPTIONS] < big_dump.sql | mysql -hmysql.foo.com -u... -p... DBNAME
```

to import a large database to mysql.foo.com while trying not to 
overload it.

## USAGE
TBD

## CAVEATS
The program simply tokenises its input by newlines. You may want to instruct 
your `mysqldump(1)` to output SQL statements line-by-line, although that 
shouldn't be strictly necessary.
