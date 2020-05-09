/*
** wikibook table contains generated table of contents
**
** uuid: same uuid as related job
** subject: starting page
** generator_version: version of wikibookgen which generated wikibook. Allows API to know if content is stale
** gen_date: date of creation
** model: model of wikibook (ABSTRACT, TOUR, ENCYCLOPEDIA)
** pages: number of pages
** table_of_content: table of content json related to model. Contains all information necessary to generate book.
*/ 
CREATE TABLE wikibook (id UUID PRIMARY KEY, subject TEXT, generator_version TEXT, gen_date TIMESTAMP WITH TIME ZONE, model TEXT, pages INT, table_of_content JSON);

/*
** job table contains wikibook orders to be generated.
**
** subject: starting page
** model: model of wikibook (ABSTRACT, TOUR, ENCYCLOPEDIA)
** creation_date: date of creation
** status: job status (CREATED, ONGOING, DONE)
*/
CREATE TABLE job (id UUID DEFAULT gen_random_uuid() PRIMARY KEY, subject TEXT, model TEXT, creation_date TIMESTAMP WITH TIME ZONE, status TEXT, book_id UUID);

