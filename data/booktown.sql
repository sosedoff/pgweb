--
-- Selected TOC Entries:
--
--
-- TOC Entry ID 1 (OID 0)
--
-- Name: booktown Type: DATABASE Owner: postgres
--

DROP DATABASE IF EXISTS "booktown";
CREATE DATABASE "booktown";

\connect booktown postgres
--
-- TOC Entry ID 2 (OID 2991542)
--
-- Name: DATABASE "booktown" Type: COMMENT Owner: 
--

COMMENT ON DATABASE "booktown" IS 'The Book Town Database.';

--
-- TOC Entry ID 33 (OID 3629264)
--
-- Name: books Type: TABLE Owner: manager
--

CREATE TABLE "books" (
	"id" integer NOT NULL,
	"title" text NOT NULL,
	"author_id" integer,
	"subject_id" integer,
	Constraint "books_id_pkey" Primary Key ("id")
);

--
-- TOC Entry ID 47 (OID 2991733)
--
-- Name: "plpgsql_call_handler" () Type: FUNCTION Owner: postgres
--

CREATE FUNCTION "plpgsql_call_handler" () RETURNS opaque AS '/usr/local/pgsql/lib/plpgsql.so', 'plpgsql_call_handler' LANGUAGE 'C';

--
-- TOC Entry ID 48 (OID 2991734)
--
-- Name: plpgsql Type: PROCEDURAL LANGUAGE Owner: 
--

CREATE TRUSTED PROCEDURAL LANGUAGE 'plpgsql' HANDLER "plpgsql_call_handler" LANCOMPILER 'PL/pgSQL';

--
-- TOC Entry ID 51 (OID 2991735)
--
-- Name: "audit_bk" (integer) Type: FUNCTION Owner: postgres
--

CREATE FUNCTION "audit_bk" (integer) RETURNS integer AS '
	DECLARE
	 key ALIAS FOR $1;
	table_data inventory%ROWTYPE;
  BEGIN
	INSERT INTO inventory_audit SELECT table_data WHERE sort_key=key;
	
	IF NOT FOUND THEN
	  RAISE EXCEPTION ''View'' || key || '' not found '';
	END IF;
	
	return 1;
end;
' LANGUAGE 'plpgsql';

--
-- TOC Entry ID 52 (OID 2991736)
--
-- Name: "audit" (integer) Type: FUNCTION Owner: postgres
--

CREATE FUNCTION "audit" (integer) RETURNS integer AS '
	DECLARE
	 key ALIAS FOR $1;
	table_data inventory%ROWTYPE;
  BEGIN
	INSERT INTO inventory_audit SELECT table_data WHERE sort_key=key;
	
	IF NOT FOUND THEN
	  RAISE EXCEPTION ''View'' || key || '' not found '';
	END IF;
	
	return 1;
end;
' LANGUAGE 'plpgsql';

--
-- TOC Entry ID 53 (OID 2991737)
--
-- Name: "auditbk" () Type: FUNCTION Owner: postgres
--

CREATE FUNCTION "auditbk" () RETURNS integer AS '
	DECLARE
	 key ALIAS FOR $1;
	table_data inventory%ROWTYPE;
  BEGIN
	INSERT INTO inventory_audit SELECT table_data WHERE sort_key=key;
	
	IF NOT FOUND THEN
	  RAISE EXCEPTION ''View'' || key || '' not found '';
	END IF;
	
	return 1;
end;
' LANGUAGE 'plpgsql';

--
-- TOC Entry ID 54 (OID 2991738)
--
-- Name: "audit_bk1" () Type: FUNCTION Owner: postgres
--

CREATE FUNCTION "audit_bk1" () RETURNS opaque AS '
	DECLARE
	 key ALIAS FOR $1;
	table_data inventory%ROWTYPE;
  BEGIN
	INSERT INTO inventory_audit SELECT table_data WHERE sort_key=key;
	
	IF NOT FOUND THEN
	  RAISE EXCEPTION ''View'' || key || '' not found '';
	END IF;
	
	return 1;
end;
' LANGUAGE 'plpgsql';

--
-- TOC Entry ID 73 (OID 2991835)
--
-- Name: "test_check_a_id" () Type: FUNCTION Owner: example
--

CREATE FUNCTION "test_check_a_id" () RETURNS opaque AS '
    BEGIN
     -- checks to make sure the author id
     -- inserted is not left blank or less than 100

        IF NEW.a_id ISNULL THEN
           RAISE EXCEPTION
           ''The author id cannot be left blank!'';
        ELSE
           IF NEW.a_id < 100 THEN
              RAISE EXCEPTION
              ''Please insert a valid author id.'';
           ELSE
           RETURN NEW;
           END IF;
        END IF;
    END;
' LANGUAGE 'plpgsql';

--
-- TOC Entry ID 66 (OID 2992619)
--
-- Name: "audit_test" () Type: FUNCTION Owner: example
--

CREATE FUNCTION "audit_test" () RETURNS opaque AS '
    BEGIN   
       
      IF TG_OP = ''INSERT'' OR TG_OP = ''UPDATE'' THEN

         NEW.user_aud := current_user;
         NEW.mod_time := ''NOW'';

        INSERT INTO inventory_audit SELECT * FROM inventory WHERE prod_id=NEW.prod_id;
              
      RETURN NEW; 

      ELSE if TG_OP = ''DELETE'' THEN
        INSERT INTO inventory_audit SELECT *, current_user, ''NOW'' FROM inventory WHERE prod_id=OLD.prod_id;

      RETURN OLD;
      END IF;
     END IF;
    END;
' LANGUAGE 'plpgsql';

--
-- TOC Entry ID 67 (OID 3000878)
--
-- Name: "first" () Type: FUNCTION Owner: example
--

CREATE FUNCTION "first" () RETURNS integer AS ' 
       DecLarE
        oNe IntEgER := 1;
       bEGiN
        ReTUrn oNE;       
       eNd;
' LANGUAGE 'plpgsql';

--
-- TOC Entry ID 68 (OID 3000881)
--
-- Name: "test" (integer) Type: FUNCTION Owner: example
--

CREATE FUNCTION "test" (integer) RETURNS integer AS '
  
 DECLARE 
   -- defines the variable as ALIAS
  variable ALIAS FOR $1;
 BEGIN
  -- displays the variable after multiplying it by two 
  return variable * 2.0;
 END; 
 ' LANGUAGE 'plpgsql';

--
-- TOC Entry ID 69 (OID 3000991)
--
-- Name: "you_me" (integer) Type: FUNCTION Owner: example
--

CREATE FUNCTION "you_me" (integer) RETURNS integer AS '
  DECLARE
   RENAME $1 TO user_no;
    --you INTEGER := 5;
  BEGIN
    return user_no;
  END;' LANGUAGE 'plpgsql';

--
-- TOC Entry ID 62 (OID 3001136)
--
-- Name: "count_by_two" (integer) Type: FUNCTION Owner: example
--

CREATE FUNCTION "count_by_two" (integer) RETURNS integer AS '
     DECLARE
          userNum ALIAS FOR $1;
          i integer;
     BEGIN
          i := 1;
          WHILE userNum[1] < 20 LOOP
                i = i+1; 
                return userNum;              
          END LOOP;
          
     END;
   ' LANGUAGE 'plpgsql';

--
-- TOC Entry ID 63 (OID 3001139)
--
-- Name: "me" () Type: FUNCTION Owner: example
--

CREATE FUNCTION "me" () RETURNS text AS '
  DECLARE
     you text := ''testing'';
     RENAME you to me;
  BEGIN
     return me;
  END;' LANGUAGE 'plpgsql';

--
-- TOC Entry ID 64 (OID 3001149)
--
-- Name: "display_cust" (integer) Type: FUNCTION Owner: example
--

CREATE FUNCTION "display_cust" (integer) RETURNS text AS '
 DECLARE
   -- declares an alias name for input
   cust_num ALIAS FOR $1;

   -- declares a row type
   cust_info customer%ROWTYPE;
 BEGIN
   -- puts information into the newly declared rowtype
   SELECT into cust_info * 
     FROM customer 
    WHERE cust_id=cust_num;    

   -- displays the customer lastname
   -- extracted from the rowtype   
   return cust_info.lastname;
 END;
 ' LANGUAGE 'plpgsql';

--
-- TOC Entry ID 65 (OID 3001151)
--
-- Name: "mixed" () Type: FUNCTION Owner: example
--

CREATE FUNCTION "mixed" () RETURNS integer AS '
       DecLarE
          --assigns 1 to the oNe variable
          oNe IntEgER 
          := 1;

       bEGiN

          --displays the value of oNe
          ReTUrn oNe;       
       eNd;
       ' LANGUAGE 'plpgsql';

--
-- TOC Entry ID 12 (OID 3117548)
--
-- Name: publishers Type: TABLE Owner: postgres
--

CREATE TABLE "publishers" (
	"id" integer NOT NULL,
	"name" text,
	"address" text,
	Constraint "publishers_pkey" Primary Key ("id")
);

--
-- TOC Entry ID 55 (OID 3117729)
--
-- Name: "compound_word" (text,text) Type: FUNCTION Owner: example
--

CREATE FUNCTION "compound_word" (text,text) RETURNS text AS '
     DECLARE
       -- defines an alias name for the two input values
       word1 ALIAS FOR $1;
       word2 ALIAS FOR $2;
     BEGIN
       -- displays the resulting joined words
       RETURN word1 || word2;
     END;
  ' LANGUAGE 'plpgsql';

--
-- TOC Entry ID 56 (OID 3117787)
--
-- Name: "givename" () Type: FUNCTION Owner: example
--

CREATE FUNCTION "givename" () RETURNS opaque AS '
 DECLARE
   tablename text;
 BEGIN
   
   tablename = TG_RELNAME; 
   INSERT INTO INVENTORY values (123, tablename);
   return old;
 END;
' LANGUAGE 'plpgsql';

--
-- TOC Entry ID 14 (OID 3389594)
--
-- Name: authors Type: TABLE Owner: manager
--

CREATE TABLE "authors" (
	"id" integer NOT NULL,
	"last_name" text,
	"first_name" text,
	Constraint "authors_pkey" Primary Key ("id")
);

--
-- TOC Entry ID 15 (OID 3389632)
--
-- Name: states Type: TABLE Owner: postgres
--

CREATE TABLE "states" (
	"id" integer NOT NULL,
	"name" text,
	"abbreviation" character(2),
	Constraint "state_pkey" Primary Key ("id")
);

--
-- TOC Entry ID 16 (OID 3389702)
--
-- Name: my_list Type: TABLE Owner: postgres
--

CREATE TABLE "my_list" (
	"todos" text
);

--
-- TOC Entry ID 17 (OID 3390348)
--
-- Name: stock Type: TABLE Owner: postgres
--

CREATE TABLE "stock" (
	"isbn" text NOT NULL,
	"cost" numeric(5,2),
	"retail" numeric(5,2),
	"stock" integer,
	Constraint "stock_pkey" Primary Key ("isbn")
);

--
-- TOC Entry ID 4 (OID 3390416)
--
-- Name: subject_ids Type: SEQUENCE Owner: postgres
--

CREATE SEQUENCE "subject_ids" start 0 increment 1 maxvalue 2147483647 minvalue 0  cache 1 ;

--
-- TOC Entry ID 19 (OID 3390653)
--
-- Name: numeric_values Type: TABLE Owner: postgres
--

CREATE TABLE "numeric_values" (
	"num" numeric(30,6)
);

--
-- TOC Entry ID 20 (OID 3390866)
--
-- Name: daily_inventory Type: TABLE Owner: postgres
--

CREATE TABLE "daily_inventory" (
	"isbn" text,
	"is_stocked" boolean
);

--
-- TOC Entry ID 21 (OID 3391084)
--
-- Name: money_example Type: TABLE Owner: postgres
--

CREATE TABLE "money_example" (
	"money_cash" money,
	"numeric_cash" numeric(6,2)
);

--
-- TOC Entry ID 22 (OID 3391184)
--
-- Name: shipments Type: TABLE Owner: postgres
--

CREATE TABLE "shipments" (
	"id" integer DEFAULT nextval('"shipments_ship_id_seq"'::text) NOT NULL,
	"customer_id" integer,
	"isbn" text,
	"ship_date" timestamp with time zone
);

--
-- TOC Entry ID 24 (OID 3391454)
--
-- Name: customers Type: TABLE Owner: manager
--

CREATE TABLE "customers" (
	"id" integer NOT NULL,
	"last_name" text,
	"first_name" text,
	Constraint "customers_pkey" Primary Key ("id")
);

--
-- TOC Entry ID 6 (OID 3574018)
--
-- Name: book_ids Type: SEQUENCE Owner: postgres
--

CREATE SEQUENCE "book_ids" start 0 increment 1 maxvalue 2147483647 minvalue 0  cache 1 ;

--
-- TOC Entry ID 26 (OID 3574043)
--
-- Name: book_queue Type: TABLE Owner: postgres
--

CREATE TABLE "book_queue" (
	"title" text NOT NULL,
	"author_id" integer,
	"subject_id" integer,
	"approved" boolean
);

--
-- TOC Entry ID 78 (OID 3574403)
--
-- Name: "title" (integer) Type: FUNCTION Owner: postgres
--

CREATE FUNCTION "title" (integer) RETURNS text AS 'SELECT title from books where id = $1' LANGUAGE 'sql';

--
-- TOC Entry ID 27 (OID 3574983)
--
-- Name: stock_backup Type: TABLE Owner: postgres
--

CREATE TABLE "stock_backup" (
	"isbn" text,
	"cost" numeric(5,2),
	"retail" numeric(5,2),
	"stock" integer
);

--
-- TOC Entry ID 89 (OID 3625934)
--
-- Name: "double_price" (double precision) Type: FUNCTION Owner: postgres
--

CREATE FUNCTION "double_price" (double precision) RETURNS double precision AS '
  DECLARE
  BEGIN
    return $1 * 2;
  END;
' LANGUAGE 'plpgsql';

--
-- TOC Entry ID 90 (OID 3625935)
--
-- Name: "triple_price" (double precision) Type: FUNCTION Owner: postgres
--

CREATE FUNCTION "triple_price" (double precision) RETURNS double precision AS '
  DECLARE
     -- Declare input_price as an alias for the
     -- argument variable normally referenced with
     -- the $1 identifier.
    input_price ALIAS FOR $1;
 
  BEGIN
     -- Return the input price multiplied by three.
    RETURN input_price * 3;
  END;
 ' LANGUAGE 'plpgsql';

--
-- TOC Entry ID 87 (OID 3625944)
--
-- Name: "stock_amount" (integer,integer) Type: FUNCTION Owner: postgres
--

CREATE FUNCTION "stock_amount" (integer,integer) RETURNS integer AS '
  DECLARE
     -- Declare aliases for function arguments.
    b_id ALIAS FOR $1;
    b_edition ALIAS FOR $2;
     -- Declare variable to store the ISBN number.
    b_isbn TEXT;
     -- Declare variable to store the stock amount.
    stock_amount INTEGER;
  BEGIN
     -- This SELECT INTO statement retrieves the ISBN
     -- number of the row in the editions table that had
     -- both the book ID number and edition number that
     -- were provided as function arguments.
    SELECT INTO b_isbn isbn FROM editions WHERE
      book_id = b_id AND edition = b_edition;
 
     -- Check to see if the ISBN number retrieved
     -- is NULL.  This will happen if there is not an
     -- existing book with both the ID number and edition
     -- number specified in the function arguments.
     -- If the ISBN is null, the function returns a
     -- value of -1 and ends.
    IF b_isbn IS NULL THEN
      RETURN -1;
    END IF;
 
     -- Retrieve the amount of books available from the
     -- stock table and record the number in the
     -- stock_amount variable.
    SELECT INTO stock_amount stock FROM stock WHERE isbn = b_isbn;
 
     -- Return the amount of books available.
    RETURN stock_amount;
  END;
' LANGUAGE 'plpgsql';

--
-- TOC Entry ID 86 (OID 3625946)
--
-- Name: "in_stock" (integer,integer) Type: FUNCTION Owner: postgres
--

CREATE FUNCTION "in_stock" (integer,integer) RETURNS boolean AS '
  DECLARE
    b_id ALIAS FOR $1;
    b_edition ALIAS FOR $2;
    b_isbn TEXT;
    stock_amount INTEGER;
  BEGIN
     -- This SELECT INTO statement retrieves the ISBN
     -- number of the row in the editions table that had
     -- both the book ID number and edition number that
     -- were provided as function arguments.
    SELECT INTO b_isbn isbn FROM editions WHERE
      book_id = b_id AND edition = b_edition;
 
     -- Check to see if the ISBN number retrieved
     -- is NULL.  This will happen if there is not an
     -- existing book with both the ID number and edition
     -- number specified in the function arguments.
     -- If the ISBN is null, the function returns a
     -- FALSE value and ends.
    IF b_isbn IS NULL THEN
      RETURN FALSE;
    END IF;
 
     -- Retrieve the amount of books available from the
     -- stock table and record the number in the
     -- stock_amount variable.
    SELECT INTO stock_amount stock FROM stock WHERE isbn = b_isbn;
 
     -- Use an IF/THEN/ELSE check to see if the amount
     -- of books available is less than, or equal to 0.
     -- If so, return FALSE.  If not, return TRUE.
    IF stock_amount <= 0 THEN
      RETURN FALSE;
    ELSE
      RETURN TRUE;
    END IF;
  END;
' LANGUAGE 'plpgsql';

--
-- TOC Entry ID 82 (OID 3626013)
--
-- Name: "extract_all_titles" () Type: FUNCTION Owner: postgres
--

CREATE FUNCTION "extract_all_titles" () RETURNS text AS '
  DECLARE
    sub_id INTEGER;
    text_output TEXT = '' '';
    sub_title TEXT;
    row_data books%ROWTYPE;
  BEGIN
    FOR i IN 0..15 LOOP
      SELECT INTO sub_title subject FROM subjects WHERE id = i;
      text_output = text_output || ''
'' || sub_title || '':
'';

      FOR row_data IN SELECT * FROM books
        WHERE subject_id = i  LOOP

        IF NOT FOUND THEN
          text_output := text_output || ''None.
'';
        ELSE
          text_output := text_output || row_data.title || ''
'';
        END IF;

      END LOOP;
    END LOOP;
    RETURN text_output;
  END;
' LANGUAGE 'plpgsql';

--
-- TOC Entry ID 79 (OID 3626052)
--
-- Name: "books_by_subject" (text) Type: FUNCTION Owner: postgres
--

CREATE FUNCTION "books_by_subject" (text) RETURNS text AS '
  DECLARE
    sub_title ALIAS FOR $1;
    sub_id INTEGER;
    found_text TEXT :='''';
  BEGIN
      SELECT INTO sub_id id FROM subjects WHERE subject = sub_title;
      RAISE NOTICE ''sub_id = %'',sub_id;
      IF sub_title = ''all'' THEN
        found_text := extract_all_titles();
        RETURN found_text;
      ELSE IF sub_id  >= 0 THEN
          found_text := extract_title(sub_id);
          RETURN  ''
'' || sub_title || '':
'' || found_text;
        END IF;
    END IF;
    RETURN ''Subject not found.'';
  END;
' LANGUAGE 'plpgsql';

--
-- TOC Entry ID 81 (OID 3626590)
--
-- Name: "add_two_loop" (integer,integer) Type: FUNCTION Owner: postgres
--

CREATE FUNCTION "add_two_loop" (integer,integer) RETURNS integer AS '
  DECLARE
 
     -- Declare aliases for function arguments.
 
    low_number ALIAS FOR $1;
    high_number ALIAS FOR $2;
 
     -- Declare a variable to hold the result.
 
    result INTEGER = 0;
 
  BEGIN
 
    WHILE result != high_number LOOP
      result := result + 1;
    END LOOP;
 
    RETURN result;
  END;
' LANGUAGE 'plpgsql';

--
-- TOC Entry ID 92 (OID 3627916)
--
-- Name: "extract_all_titles2" () Type: FUNCTION Owner: postgres
--

CREATE FUNCTION "extract_all_titles2" () RETURNS text AS '
  DECLARE
    sub_id INTEGER;
    text_output TEXT = '' '';
    sub_title TEXT;
    row_data books%ROWTYPE;
  BEGIN
    FOR i IN 0..15 LOOP
      SELECT INTO sub_title subject FROM subjects WHERE id = i;
      text_output = text_output || ''
'' || sub_title || '':
'';

      FOR row_data IN SELECT * FROM books
        WHERE subject_id = i  LOOP

        text_output := text_output || row_data.title || ''
'';

      END LOOP;
    END LOOP;
    RETURN text_output;
  END;
' LANGUAGE 'plpgsql';

--
-- TOC Entry ID 94 (OID 3627974)
--
-- Name: "extract_title" (integer) Type: FUNCTION Owner: postgres
--

CREATE FUNCTION "extract_title" (integer) RETURNS text AS '
  DECLARE
    sub_id ALIAS FOR $1;
    text_output TEXT :=''
'';
    row_data RECORD;
  BEGIN
    FOR row_data IN SELECT * FROM books
    WHERE subject_id = sub_id ORDER BY title  LOOP
      text_output := text_output || row_data.title || ''
'';
    END LOOP;
    RETURN text_output;
  END;
' LANGUAGE 'plpgsql';

--
-- TOC Entry ID 95 (OID 3628021)
--
-- Name: "raise_test" () Type: FUNCTION Owner: postgres
--

CREATE FUNCTION "raise_test" () RETURNS integer AS '
  DECLARE
 
     -- Declare an integer variable for testing.
 
    an_integer INTEGER = 1;
 
  BEGIN
 
     -- Raise a debug level message.
 
    RAISE DEBUG ''The raise_test() function began.'';
 
    an_integer = an_integer + 1;
 
     -- Raise a notice stating that the an_integer
     -- variable was changed, then raise another notice
     -- stating its new value.
 
    RAISE NOTICE ''Variable an_integer was changed.'';
    RAISE NOTICE ''Variable an_integer value is now %.'',an_integer;
 
     -- Raise an exception.
 
    RAISE EXCEPTION ''Variable % changed.  Aborting transaction.'',an_integer;
 
  END;
' LANGUAGE 'plpgsql';

--
-- TOC Entry ID 93 (OID 3628069)
--
-- Name: "add_shipment" (integer,text) Type: FUNCTION Owner: postgres
--

CREATE FUNCTION "add_shipment" (integer,text) RETURNS timestamp with time zone AS '
  DECLARE
    customer_id ALIAS FOR $1;
    isbn ALIAS FOR $2;
    shipment_id INTEGER;
    right_now timestamp;
  BEGIN
    right_now := ''now'';
    SELECT INTO shipment_id id FROM shipments ORDER BY id DESC;
    shipment_id := shipment_id + 1;
    INSERT INTO shipments VALUES ( shipment_id, customer_id, isbn, right_now );
    RETURN right_now;
  END;
' LANGUAGE 'plpgsql';

--
-- TOC Entry ID 102 (OID 3628076)
--
-- Name: "ship_item" (text,text,text) Type: FUNCTION Owner: postgres
--

CREATE FUNCTION "ship_item" (text,text,text) RETURNS integer AS '
  DECLARE
    l_name ALIAS FOR $1;
    f_name ALIAS FOR $2;
    book_isbn ALIAS FOR $3;
    book_id INTEGER;
    customer_id INTEGER;
 
  BEGIN
 
    SELECT INTO customer_id get_customer_id(l_name,f_name);
 
    IF customer_id = -1 THEN
      RETURN -1;
    END IF;
 
    SELECT INTO book_id book_id FROM editions WHERE isbn = book_isbn;
 
    IF NOT FOUND THEN
      RETURN -1;
    END IF;
 
    PERFORM add_shipment(customer_id,book_isbn);
 
    RETURN 1;
  END;
' LANGUAGE 'plpgsql';

--
-- TOC Entry ID 103 (OID 3628114)
--
-- Name: "check_book_addition" () Type: FUNCTION Owner: postgres
--

CREATE FUNCTION "check_book_addition" () RETURNS opaque AS '
  DECLARE 
    id_number INTEGER;
    book_isbn TEXT;
  BEGIN

    SELECT INTO id_number id FROM customers WHERE id = NEW.customer_id; 

    IF NOT FOUND THEN
      RAISE EXCEPTION ''Invalid customer ID number.'';  
    END IF;

    SELECT INTO book_isbn isbn FROM editions WHERE isbn = NEW.isbn; 

    IF NOT FOUND THEN
      RAISE EXCEPTION ''Invalid ISBN.''; 
    END IF; 

    UPDATE stock SET stock = stock -1 WHERE isbn = NEW.isbn; 

    RETURN NEW; 
  END;
' LANGUAGE 'plpgsql';

--
-- TOC Entry ID 28 (OID 3628246)
--
-- Name: stock_view Type: VIEW Owner: postgres
--

CREATE VIEW "stock_view" as SELECT stock.isbn, stock.retail, stock.stock FROM stock;

CREATE MATERIALIZED VIEW "m_stock_view" as SELECT stock.isbn, stock.retail, stock.stock FROM stock;

--
-- TOC Entry ID 30 (OID 3628247)
--
-- Name: favorite_books Type: TABLE Owner: manager
--

CREATE TABLE "favorite_books" (
	"employee_id" integer,
	"books" text[]
);

--
-- TOC Entry ID 8 (OID 3628626)
--
-- Name: shipments_ship_id_seq Type: SEQUENCE Owner: manager
--

CREATE SEQUENCE "shipments_ship_id_seq" start 0 increment 1 maxvalue 2147483647 minvalue 0  cache 1 ;

--
-- TOC Entry ID 74 (OID 3628648)
--
-- Name: "check_shipment_addition" () Type: FUNCTION Owner: postgres
--

CREATE FUNCTION "check_shipment_addition" () RETURNS opaque AS '
  DECLARE
     -- Declare a variable to hold the customer ID.
    id_number INTEGER;
 
     -- Declare a variable to hold the ISBN.
    book_isbn TEXT;
  BEGIN
 
     -- If there is an ID number that matches the customer ID in
     -- the new table, retrieve it from the customers table.
    SELECT INTO id_number id FROM customers WHERE id = NEW.customer_id;
 
     -- If there was no matching ID number, raise an exception.
    IF NOT FOUND THEN
      RAISE EXCEPTION ''Invalid customer ID number.'';
    END IF;
 
     -- If there is an ISBN that matches the ISBN specified in the
     -- new table, retrieve it from the editions table.
    SELECT INTO book_isbn isbn FROM editions WHERE isbn = NEW.isbn;
 
     -- If there is no matching ISBN, raise an exception.
    IF NOT FOUND THEN
      RAISE EXCEPTION ''Invalid ISBN.'';
    END IF;
 
    -- If the previous checks succeeded, update the stock amount
    -- for INSERT commands.
    IF TG_OP = ''INSERT'' THEN
       UPDATE stock SET stock = stock -1 WHERE isbn = NEW.isbn;
    END IF;
 
    RETURN NEW;
  END;
' LANGUAGE 'plpgsql';

--
-- TOC Entry ID 31 (OID 3628899)
--
-- Name: employees Type: TABLE Owner: postgres
--

CREATE TABLE "employees" (
	"id" integer NOT NULL,
	"last_name" text NOT NULL,
	"first_name" text,
	CONSTRAINT "employees_id" CHECK ((id > 100)),
	Constraint "employees_pkey" Primary Key ("id")
);

--
-- TOC Entry ID 32 (OID 3629174)
--
-- Name: editions Type: TABLE Owner: manager
--

CREATE TABLE "editions" (
	"isbn" text NOT NULL,
	"book_id" integer,
	"edition" integer,
	"publisher_id" integer,
	"publication" date,
	"type" character(1),
	CONSTRAINT "integrity" CHECK (((book_id NOTNULL) AND (edition NOTNULL))),
	Constraint "pkey" Primary Key ("isbn")
);

--
-- TOC Entry ID 10 (OID 3629402)
--
-- Name: author_ids Type: SEQUENCE Owner: manager
--

CREATE SEQUENCE "author_ids" start 0 increment 1 maxvalue 2147483647 minvalue 0  cache 1 ;

--
-- TOC Entry ID 35 (OID 3629424)
--
-- Name: distinguished_authors Type: TABLE Owner: manager
--

CREATE TABLE "distinguished_authors" (
	"award" text
)
INHERITS ("authors");

--
-- TOC Entry ID 107 (OID 3726476)
--
-- Name: "isbn_to_title" (text) Type: FUNCTION Owner: manager
--

CREATE FUNCTION "isbn_to_title" (text) RETURNS text AS 'SELECT title FROM books
                                 JOIN editions AS e (isbn, id)
                                 USING (id)
                                 WHERE isbn = $1' LANGUAGE 'sql';

--
-- TOC Entry ID 36 (OID 3727889)
--
-- Name: favorite_authors Type: TABLE Owner: manager
--

CREATE TABLE "favorite_authors" (
	"employee_id" integer,
	"authors_and_titles" text[]
);

--
-- TOC Entry ID 99 (OID 3728728)
--
-- Name: "get_customer_name" (integer) Type: FUNCTION Owner: postgres
--

CREATE FUNCTION "get_customer_name" (integer) RETURNS text AS '
  DECLARE
  
    -- Declare aliases for user input.
    customer_id ALIAS FOR $1;
    
    -- Declare variables to hold the customer name.
    customer_fname TEXT;
    customer_lname TEXT;
  
  BEGIN
  
    -- Retrieve the customer first and last name for the customer whose
    -- ID matches the value supplied as a function argument.
    SELECT INTO customer_fname, customer_lname 
                first_name, last_name FROM customers
      WHERE id = customer_id;
    
    -- Return the name.
    RETURN customer_fname || '' '' || customer_lname;
  END;
' LANGUAGE 'plpgsql';

--
-- TOC Entry ID 100 (OID 3728729)
--
-- Name: "get_customer_id" (text,text) Type: FUNCTION Owner: postgres
--

CREATE FUNCTION "get_customer_id" (text,text) RETURNS integer AS '
  DECLARE
 
    -- Declare aliases for user input.
    l_name ALIAS FOR $1;
    f_name ALIAS FOR $2;
 
    -- Declare a variable to hold the customer ID number.
    customer_id INTEGER;
 
  BEGIN
 
    -- Retrieve the customer ID number of the customer whose first and last
    --  name match the values supplied as function arguments.
    SELECT INTO customer_id id FROM customers
      WHERE last_name = l_name AND first_name = f_name;
 
    -- Return the ID number.
    RETURN customer_id;
  END;
' LANGUAGE 'plpgsql';

--
-- TOC Entry ID 101 (OID 3728730)
--
-- Name: "get_author" (text) Type: FUNCTION Owner: postgres
--

CREATE FUNCTION "get_author" (text) RETURNS text AS '
  DECLARE
 
      -- Declare an alias for the function argument,
      -- which should be the first name of an author.
     f_name ALIAS FOR $1;
 
       -- Declare a variable with the same type as
       -- the last_name field of the authors table.
     l_name authors.last_name%TYPE;
 
  BEGIN
 
      -- Retrieve the last name of an author from the
      -- authors table whose first name matches the
      -- argument received by the function, and
      -- insert it into the l_name variable.
     SELECT INTO l_name last_name FROM authors WHERE first_name = f_name;
 
       -- Return the first name and last name, separated
       -- by a space.
     return f_name || '' '' || l_name;
 
  END;
' LANGUAGE 'plpgsql';

--
-- TOC Entry ID 97 (OID 3728759)
--
-- Name: "get_author" (integer) Type: FUNCTION Owner: postgres
--

CREATE FUNCTION "get_author" (integer) RETURNS text AS '
  DECLARE
 
    -- Declare an alias for the function argument,
    -- which should be the id of the author.
    author_id ALIAS FOR $1;
 
    -- Declare a variable that uses the structure of
    -- the authors table.
    found_author authors%ROWTYPE;
 
  BEGIN
 
    -- Retrieve a row of author information for
    -- the author whose id number matches
    -- the argument received by the function.
    SELECT INTO found_author * FROM authors WHERE id = author_id;
 
    -- Return the first
    RETURN found_author.first_name || '' '' || found_author.last_name;
 
  END;
' LANGUAGE 'plpgsql';

--
-- TOC Entry ID 70 (OID 3743412)
--
-- Name: "html_linebreaks" (text) Type: FUNCTION Owner: postgres
--

CREATE FUNCTION "html_linebreaks" (text) RETURNS text AS '
  DECLARE
    formatted_string text := '''';
  BEGIN
    FOR i IN 0 .. length($1) LOOP
      IF substr($1, i, 1) = ''
'' THEN
        formatted_string := formatted_string || ''<br>'';
      ELSE
        formatted_string := formatted_string || substr($1, i, 1);
      END IF;
    END LOOP;
    RETURN formatted_string;
  END;
' LANGUAGE 'plpgsql';

--
-- TOC Entry ID 37 (OID 3751599)
--
-- Name: text_sorting Type: TABLE Owner: postgres
--

CREATE TABLE "text_sorting" (
	"letter" character(1)
);

--
-- TOC Entry ID 38 (OID 3751882)
--
-- Name: subjects Type: TABLE Owner: postgres
--

CREATE TABLE "subjects" (
	"id" integer NOT NULL,
	"subject" text,
	"location" text,
	Constraint "subjects_pkey" Primary Key ("id")
);

--
-- TOC Entry ID 108 (OID 3751924)
--
-- Name: sum(text) Type: AGGREGATE Owner: postgres
--

CREATE AGGREGATE sum ( BASETYPE = text, SFUNC = textcat, STYPE = text, INITCOND = '' );

--
-- TOC Entry ID 39 (OID 3751975)
--
-- Name: alternate_stock Type: TABLE Owner: postgres
--

CREATE TABLE "alternate_stock" (
	"isbn" text,
	"cost" numeric(5,2),
	"retail" numeric(5,2),
	"stock" integer
);

--
-- TOC Entry ID 40 (OID 3752020)
--
-- Name: book_backup Type: TABLE Owner: postgres
--

CREATE TABLE "book_backup" (
	"id" integer,
	"title" text,
	"author_id" integer,
	"subject_id" integer
);

--
-- TOC Entry ID 80 (OID 3752102)
--
-- Name: "sync_authors_and_books" () Type: FUNCTION Owner: postgres
--

CREATE FUNCTION "sync_authors_and_books" () RETURNS opaque AS '
  BEGIN
    IF TG_OP = ''UPDATE'' THEN
      UPDATE books SET author_id = new.id WHERE author_id = old.id; 
    END IF;
    RETURN new;
  END;
' LANGUAGE 'plpgsql';

--
-- TOC Entry ID 41 (OID 4063343)
--
-- Name: schedules Type: TABLE Owner: postgres
--
 
CREATE TABLE "schedules" (
        "employee_id" integer NOT NULL,
        "schedule" text,
        Constraint "schedules_pkey" Primary Key ("employee_id")
);

--
-- TOC Entry ID 42 (OID 4063653)
--
-- Name: recent_shipments Type: VIEW Owner: postgres
--

CREATE VIEW "recent_shipments" as SELECT count(*) AS num_shipped, max(shipments.ship_date) AS max, b.title FROM ((shipments JOIN editions USING (isbn)) NATURAL JOIN books b(book_id)) GROUP BY b.title ORDER BY count(*) DESC;

--
-- Data for TOC Entry ID 112 (OID 3117548)
--
-- Name: publishers Type: TABLE DATA Owner: postgres
--


COPY "publishers"  FROM stdin;
150	Kids Can Press	Kids Can Press, 29 Birch Ave. Toronto, ON  M4V 1E2
91	Henry Holt & Company, Inc.	Henry Holt & Company, Inc. 115 West 18th Street New York, NY 10011
113	O'Reilly & Associates	O'Reilly & Associates, Inc. 101 Morris St, Sebastopol, CA 95472
62	Watson-Guptill Publications	1515 Boradway, New York, NY 10036
105	Noonday Press	Farrar Straus & Giroux Inc, 19 Union Square W, New York, NY 10003
99	Ace Books	The Berkley Publishing Group, Penguin Putnam Inc, 375 Hudson St, New York, NY 10014
101	Roc	Penguin Putnam Inc, 375 Hudson St, New York, NY 10014
163	Mojo Press	Mojo Press, PO Box 1215, Dripping Springs, TX 78720
171	Books of Wonder	Books of Wonder, 16 W. 18th St. New York, NY, 10011
102	Penguin	Penguin Putnam Inc, 375 Hudson St, New York, NY 10014
75	Doubleday	Random House, Inc, 1540 Broadway, New York, NY 10036
65	HarperCollins	HarperCollins Publishers, 10 E 53rd St, New York, NY 10022
59	Random House	Random House, Inc, 1540 Broadway, New York, NY 10036
\.
--
-- Data for TOC Entry ID 113 (OID 3389594)
--
-- Name: authors Type: TABLE DATA Owner: manager
--


COPY "authors"  FROM stdin;
1111	Denham	Ariel
1212	Worsley	John
15990	Bourgeois	Paulette
25041	Bianco	Margery Williams
16	Alcott	Louisa May
4156	King	Stephen
1866	Herbert	Frank
1644	Hogarth	Burne
2031	Brown	Margaret Wise
115	Poe	Edgar Allen
7805	Lutz	Mark
7806	Christiansen	Tom
1533	Brautigan	Richard
1717	Brite	Poppy Z.
2112	Gorey	Edward
2001	Clarke	Arthur C.
1213	Brookins	Andrew
\.
--
-- Data for TOC Entry ID 114 (OID 3389632)
--
-- Name: states Type: TABLE DATA Owner: postgres
--


COPY "states"  FROM stdin;
42	Washington	WA
51	Oregon	OR
\.
--
-- Data for TOC Entry ID 115 (OID 3389702)
--
-- Name: my_list Type: TABLE DATA Owner: postgres
--


COPY "my_list"  FROM stdin;
Pick up laundry.
Send out bills.
Wrap up Grand Unifying Theory for publication.
\.
--
-- Data for TOC Entry ID 116 (OID 3390348)
--
-- Name: stock Type: TABLE DATA Owner: postgres
--


COPY "stock"  FROM stdin;
0385121679	29.00	36.95	65
039480001X	30.00	32.95	31
0394900014	23.00	23.95	0
044100590X	36.00	45.95	89
0441172717	17.00	21.95	77
0451160916	24.00	28.95	22
0451198492	36.00	46.95	0
0451457994	17.00	22.95	0
0590445065	23.00	23.95	10
0679803335	20.00	24.95	18
0694003611	25.00	28.95	50
0760720002	18.00	23.95	28
0823015505	26.00	28.95	16
0929605942	19.00	21.95	25
1885418035	23.00	24.95	77
0394800753	16.00	16.95	4
\.
--
-- Data for TOC Entry ID 117 (OID 3390653)
--
-- Name: numeric_values Type: TABLE DATA Owner: postgres
--


COPY "numeric_values"  FROM stdin;
68719476736.000000
68719476737.000000
6871947673778.000000
999999999999999999999999.999900
999999999999999999999999.999999
-999999999999999999999999.999999
-100000000000000000000000.999999
1.999999
2.000000
2.000000
999999999999999999999999.999999
999999999999999999999999.000000
\.
--
-- Data for TOC Entry ID 118 (OID 3390866)
--
-- Name: daily_inventory Type: TABLE DATA Owner: postgres
--


COPY "daily_inventory"  FROM stdin;
039480001X	t
044100590X	t
0451198492	f
0394900014	f
0441172717	t
0451160916	f
0385121679	\N
\.
--
-- Data for TOC Entry ID 119 (OID 3391084)
--
-- Name: money_example Type: TABLE DATA Owner: postgres
--


COPY "money_example"  FROM stdin;
$12.24	12.24
\.
--
-- Data for TOC Entry ID 120 (OID 3391184)
--
-- Name: shipments Type: TABLE DATA Owner: postgres
--


COPY "shipments"  FROM stdin;
375	142	039480001X	2001-08-06 09:29:21-07
323	671	0451160916	2001-08-14 10:36:41-07
998	1045	0590445065	2001-08-12 12:09:47-07
749	172	0694003611	2001-08-11 10:52:34-07
662	655	0679803335	2001-08-09 07:30:07-07
806	1125	0760720002	2001-08-05 09:34:04-07
102	146	0394900014	2001-08-11 13:34:08-07
813	112	0385121679	2001-08-08 09:53:46-07
652	724	1885418035	2001-08-14 13:41:39-07
599	430	0929605942	2001-08-10 08:29:42-07
969	488	0441172717	2001-08-14 08:42:58-07
433	898	044100590X	2001-08-12 08:46:35-07
660	409	0451457994	2001-08-07 11:56:42-07
310	738	0451198492	2001-08-15 14:02:01-07
510	860	0823015505	2001-08-14 07:33:47-07
997	185	039480001X	2001-08-10 13:47:52-07
999	221	0451160916	2001-08-14 13:45:51-07
56	880	0590445065	2001-08-14 13:49:00-07
72	574	0694003611	2001-08-06 07:49:44-07
146	270	039480001X	2001-08-13 09:42:10-07
981	652	0451160916	2001-08-08 08:36:44-07
95	480	0590445065	2001-08-10 07:29:52-07
593	476	0694003611	2001-08-15 11:57:40-07
977	853	0679803335	2001-08-09 09:30:46-07
117	185	0760720002	2001-08-07 13:00:48-07
406	1123	0394900014	2001-08-13 09:47:04-07
340	1149	0385121679	2001-08-12 13:39:22-07
871	388	1885418035	2001-08-07 11:31:57-07
1000	221	039480001X	2001-09-14 16:46:32-07
1001	107	039480001X	2001-09-14 17:42:22-07
754	107	0394800753	2001-08-11 09:55:05-07
458	107	0394800753	2001-08-07 10:58:36-07
189	107	0394800753	2001-08-06 11:46:36-07
720	107	0394800753	2001-08-08 10:46:13-07
1002	107	0394800753	2001-09-22 11:23:28-07
2	107	0394800753	2001-09-22 20:58:56-07
\.
--
-- Data for TOC Entry ID 121 (OID 3391454)
--
-- Name: customers Type: TABLE DATA Owner: manager
--


COPY "customers"  FROM stdin;
107	Jackson	Annie
112	Gould	Ed
142	Allen	Chad
146	Williams	James
172	Brown	Richard
185	Morrill	Eric
221	King	Jenny
270	Bollman	Julie
388	Morrill	Royce
409	Holloway	Christine
430	Black	Jean
476	Clark	James
480	Thomas	Rich
488	Young	Trevor
574	Bennett	Laura
652	Anderson	Jonathan
655	Olson	Dave
671	Brown	Chuck
723	Eisele	Don
724	Holloway	Adam
738	Gould	Shirley
830	Robertson	Royce
853	Black	Wendy
860	Owens	Tim
880	Robinson	Tammy
898	Gerdes	Kate
964	Gould	Ramon
1045	Owens	Jean
1125	Bollman	Owen
1149	Becker	Owen
1123	Corner	Kathy
\.
--
-- Data for TOC Entry ID 122 (OID 3574043)
--
-- Name: book_queue Type: TABLE DATA Owner: postgres
--


COPY "book_queue"  FROM stdin;
Learning Python	7805	4	t
Perl Cookbook	7806	4	t
\.
--
-- Data for TOC Entry ID 123 (OID 3574983)
--
-- Name: stock_backup Type: TABLE DATA Owner: postgres
--


COPY "stock_backup"  FROM stdin;
0385121679	29.00	36.95	65
039480001X	30.00	32.95	31
0394800753	16.00	16.95	0
0394900014	23.00	23.95	0
044100590X	36.00	45.95	89
0441172717	17.00	21.95	77
0451160916	24.00	28.95	22
0451198492	36.00	46.95	0
0451457994	17.00	22.95	0
0590445065	23.00	23.95	10
0679803335	20.00	24.95	18
0694003611	25.00	28.95	50
0760720002	18.00	23.95	28
0823015505	26.00	28.95	16
0929605942	19.00	21.95	25
1885418035	23.00	24.95	77
\.
--
-- Data for TOC Entry ID 124 (OID 3628247)
--
-- Name: favorite_books Type: TABLE DATA Owner: manager
--


COPY "favorite_books"  FROM stdin;
102	{"The Hitchhiker's Guide to the Galaxy","The Restauraunt at the End of the Universe"}
103	{"There and Back Again: A Hobbit's Holiday","Kittens Squared"}
\.
--
-- Data for TOC Entry ID 125 (OID 3628899)
--
-- Name: employees Type: TABLE DATA Owner: postgres
--


COPY "employees"  FROM stdin;
101	Appel	Vincent
102	Holloway	Michael
105	Connoly	Sarah
104	Noble	Ben
103	Joble	David
106	Hall	Timothy
1008	Williams	\N
\.
--
-- Data for TOC Entry ID 126 (OID 3629174)
--
-- Name: editions Type: TABLE DATA Owner: manager
--


COPY "editions"  FROM stdin;
039480001X	1608	1	59	1957-03-01	h
0451160916	7808	1	75	1981-08-01	p
0394800753	1590	1	59	1949-03-01	p
0590445065	25908	1	150	1987-03-01	p
0694003611	1501	1	65	1947-03-04	p
0679803335	1234	1	102	1922-01-01	p
0760720002	190	1	91	1868-01-01	p
0394900014	1608	1	59	1957-01-01	p
0385121679	7808	2	75	1993-10-01	h
1885418035	156	1	163	1995-03-28	p
0929605942	156	2	171	1998-12-01	p
0441172717	4513	2	99	1998-09-01	p
044100590X	4513	3	99	1999-10-01	h
0451457994	4267	3	101	2000-09-12	p
0451198492	4267	3	101	1999-10-01	h
0823015505	2038	1	62	1958-01-01	p
0596000855	41473	2	113	2001-03-01	p
\.
--
-- Data for TOC Entry ID 127 (OID 3629264)
--
-- Name: books Type: TABLE DATA Owner: manager
--


COPY "books"  FROM stdin;
7808	The Shining	4156	9
4513	Dune	1866	15
4267	2001: A Space Odyssey	2001	15
1608	The Cat in the Hat	1809	2
1590	Bartholomew and the Oobleck	1809	2
25908	Franklin in the Dark	15990	2
1501	Goodnight Moon	2031	2
190	Little Women	16	6
1234	The Velveteen Rabbit	25041	3
2038	Dynamic Anatomy	1644	0
156	The Tell-Tale Heart	115	9
41473	Programming Python	7805	4
41477	Learning Python	7805	4
41478	Perl Cookbook	7806	4
41472	Practical PostgreSQL	1212	4
\.
--
-- Data for TOC Entry ID 128 (OID 3629424)
--
-- Name: distinguished_authors Type: TABLE DATA Owner: manager
--


COPY "distinguished_authors"  FROM stdin;
25043	Simon	Neil	Pulitzer Prize
1809	Geisel	Theodor Seuss	Pulitzer Prize
\.
--
-- Data for TOC Entry ID 129 (OID 3727889)
--
-- Name: favorite_authors Type: TABLE DATA Owner: manager
--


COPY "favorite_authors"  FROM stdin;
102	{{"J.R.R. Tolkien","The Silmarillion"},{"Charles Dickens","Great Expectations"},{"Ariel Denham","Attic Lives"}}
\.
--
-- Data for TOC Entry ID 130 (OID 3751599)
--
-- Name: text_sorting Type: TABLE DATA Owner: postgres
--


COPY "text_sorting"  FROM stdin;
0
1
2
3
A
B
C
D
a
b
c
d
\.
--
-- Data for TOC Entry ID 131 (OID 3751882)
--
-- Name: subjects Type: TABLE DATA Owner: postgres
--


COPY "subjects"  FROM stdin;
0	Arts	Creativity St
1	Business	Productivity Ave
2	Children's Books	Kids Ct
3	Classics	Academic Rd
4	Computers	Productivity Ave
5	Cooking	Creativity St
6	Drama	Main St
7	Entertainment	Main St
8	History	Academic Rd
9	Horror	Black Raven Dr
10	Mystery	Black Raven Dr
11	Poetry	Sunset Dr
12	Religion	\N
13	Romance	Main St
14	Science	Productivity Ave
15	Science Fiction	Main St
\.
--
-- Data for TOC Entry ID 132 (OID 3751975)
--
-- Name: alternate_stock Type: TABLE DATA Owner: postgres
--


COPY "alternate_stock"  FROM stdin;
0385121679	29.00	36.95	65
039480001X	30.00	32.95	31
0394900014	23.00	23.95	0
044100590X	36.00	45.95	89
0441172717	17.00	21.95	77
0451160916	24.00	28.95	22
0451198492	36.00	46.95	0
0451457994	17.00	22.95	0
0590445065	23.00	23.95	10
0679803335	20.00	24.95	18
0694003611	25.00	28.95	50
0760720002	18.00	23.95	28
0823015505	26.00	28.95	16
0929605942	19.00	21.95	25
1885418035	23.00	24.95	77
0394800753	16.00	16.95	4
\.
--
-- Data for TOC Entry ID 133 (OID 3752020)
--
-- Name: book_backup Type: TABLE DATA Owner: postgres
--


COPY "book_backup"  FROM stdin;
7808	The Shining	4156	9
4513	Dune	1866	15
4267	2001: A Space Odyssey	2001	15
1608	The Cat in the Hat	1809	2
1590	Bartholomew and the Oobleck	1809	2
25908	Franklin in the Dark	15990	2
1501	Goodnight Moon	2031	2
190	Little Women	16	6
1234	The Velveteen Rabbit	25041	3
2038	Dynamic Anatomy	1644	0
156	The Tell-Tale Heart	115	9
41472	Practical PostgreSQL	1212	4
41473	Programming Python	7805	4
41477	Learning Python	7805	4
41478	Perl Cookbook	7806	4
7808	The Shining	4156	9
4513	Dune	1866	15
4267	2001: A Space Odyssey	2001	15
1608	The Cat in the Hat	1809	2
1590	Bartholomew and the Oobleck	1809	2
25908	Franklin in the Dark	15990	2
1501	Goodnight Moon	2031	2
190	Little Women	16	6
1234	The Velveteen Rabbit	25041	3
2038	Dynamic Anatomy	1644	0
156	The Tell-Tale Heart	115	9
41473	Programming Python	7805	4
41477	Learning Python	7805	4
41478	Perl Cookbook	7806	4
41472	Practical PostgreSQL	1212	4
\.
--
-- Data for TOC Entry ID 134 (OID 4063343)
--
-- Name: schedules Type: TABLE DATA Owner: postgres
--


COPY "schedules"  FROM stdin;
102	Mon - Fri, 9am - 5pm
\.
--
-- TOC Entry ID 45 (OID 3117548)
--
-- Name: "unique_publisher_idx" Type: INDEX Owner: postgres
--

CREATE UNIQUE INDEX "unique_publisher_idx" on "publishers" using btree ( "name" "text_ops" );

--
-- TOC Entry ID 43 (OID 3391184)
--
-- Name: "shipments_ship_id_key" Type: INDEX Owner: postgres
--

CREATE UNIQUE INDEX "shipments_ship_id_key" on "shipments" using btree ( "id" "int4_ops" );

--
-- TOC Entry ID 44 (OID 3629264)
--
-- Name: "books_title_idx" Type: INDEX Owner: manager
--

CREATE  INDEX "books_title_idx" on "books" using btree ( "title" "text_ops" );

--
-- TOC Entry ID 46 (OID 3751599)
--
-- Name: "text_idx" Type: INDEX Owner: postgres
--

CREATE  INDEX "text_idx" on "text_sorting" using btree ( "letter" "bpchar_ops" );

--
-- TOC Entry ID 136 (OID 3628649)
--
-- Name: check_shipment Type: TRIGGER Owner: postgres
--

CREATE TRIGGER "check_shipment" BEFORE INSERT OR UPDATE ON "shipments"  FOR EACH ROW EXECUTE PROCEDURE "check_shipment_addition" ();

--
-- TOC Entry ID 135 (OID 3752103)
--
-- Name: sync_authors_books Type: TRIGGER Owner: manager
--

CREATE TRIGGER "sync_authors_books" BEFORE UPDATE ON "authors"  FOR EACH ROW EXECUTE PROCEDURE "sync_authors_and_books" ();

--
-- TOC Entry ID 139 (OID 4063374)
--
-- Name: "RI_ConstraintTrigger_4063373" Type: TRIGGER Owner: postgres
--

CREATE CONSTRAINT TRIGGER "valid_employee" AFTER INSERT OR UPDATE ON "schedules"  FROM "employees" NOT DEFERRABLE INITIALLY IMMEDIATE FOR EACH ROW EXECUTE PROCEDURE "RI_FKey_check_ins" ('valid_employee', 'schedules', 'employees', 'FULL', 'employee_id', 'id');

--
-- TOC Entry ID 137 (OID 4063376)
--
-- Name: "RI_ConstraintTrigger_4063375" Type: TRIGGER Owner: postgres
--

CREATE CONSTRAINT TRIGGER "valid_employee" AFTER DELETE ON "employees"  FROM "schedules" NOT DEFERRABLE INITIALLY IMMEDIATE FOR EACH ROW EXECUTE PROCEDURE "RI_FKey_noaction_del" ('valid_employee', 'schedules', 'employees', 'FULL', 'employee_id', 'id');

--
-- TOC Entry ID 138 (OID 4063378)
--
-- Name: "RI_ConstraintTrigger_4063377" Type: TRIGGER Owner: postgres
--

CREATE CONSTRAINT TRIGGER "valid_employee" AFTER UPDATE ON "employees"  FROM "schedules" NOT DEFERRABLE INITIALLY IMMEDIATE FOR EACH ROW EXECUTE PROCEDURE "RI_FKey_noaction_upd" ('valid_employee', 'schedules', 'employees', 'FULL', 'employee_id', 'id');

--
-- TOC Entry ID 140 (OID 3752079)
--
-- Name: sync_stock_with_editions Type: RULE Owner: manager
--

CREATE RULE sync_stock_with_editions AS ON UPDATE TO editions DO UPDATE stock SET isbn = new.isbn WHERE (stock.isbn = old.isbn);
--
-- TOC Entry ID 5 (OID 3390416)
--
-- Name: subject_ids Type: SEQUENCE SET Owner: 
--

SELECT setval ('"subject_ids"', 15, 't');

--
-- TOC Entry ID 7 (OID 3574018)
--
-- Name: book_ids Type: SEQUENCE SET Owner: 
--

SELECT setval ('"book_ids"', 41478, 't');

--
-- TOC Entry ID 9 (OID 3628626)
--
-- Name: shipments_ship_id_seq Type: SEQUENCE SET Owner: 
--

SELECT setval ('"shipments_ship_id_seq"', 1011, 't');

--
-- TOC Entry ID 11 (OID 3629402)
--
-- Name: author_ids Type: SEQUENCE SET Owner: 
--

SELECT setval ('"author_ids"', 25044, 't');

