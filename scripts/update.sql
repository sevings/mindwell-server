CREATE INDEX "index_tag_search" ON "mindwell"."tags" USING GIST("tag" gist_trgm_ops);
