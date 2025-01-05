CREATE TABLE liked_flavours (
  flavour_id UUID REFERENCES flavours (id) ON DELETE CASCADE NOT NULL,
  user_id UUID REFERENCES users (id) ON DELETE CASCADE NOT NULL,

  PRIMARY KEY(flavour_id, user_id)
);
