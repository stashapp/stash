CREATE EXTENSION IF NOT EXISTS plperl;
CREATE OR REPLACE FUNCTION regex_match(text, text)
RETURNS boolean AS $$
    my ($string, $pattern) = @_;
    return $string =~ /$pattern/u ? 1 : 0;
$$ LANGUAGE plperl IMMUTABLE;
