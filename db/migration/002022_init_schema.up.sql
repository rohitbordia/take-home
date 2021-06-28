CREATE TABLE script (
                        id SERIAL PRIMARY KEY,
                        name    varchar(40),
                        status   varchar(40),
                        content varchar(400),
                        last_run_status varchar(40),
                        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);