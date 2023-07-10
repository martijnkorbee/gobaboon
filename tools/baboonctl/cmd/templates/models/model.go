package models

import (
    up "github.com/upper/db/v4"
    "time"
)

// $MODELNAME$ struct
type $MODELNAME$ struct {
    ID        int         `db:"id,omitempty"`
    CreatedAt interface{} `db:"created_at"`
    UpdatedAt interface{} `db:"updated_at"`
}

// Table returns the table name
func (t *$MODELNAME$) Table() string {
    return "$TABLENAME$"
}

// Insert inserts a model into the database, using upper
func (t *$MODELNAME$) Add$MODELNAME$(m $MODELNAME$) (int, error) {
    m.CreatedAt = time.Now().String()
    m.UpdatedAt = time.Now().String()
    collection := upper.Collection(t.Table())
    res, err := collection.Insert(m)
    if err != nil {
        return 0, err
    }

    id := getInsertID(res.ID())

    return id, nil
}

// GetAll gets all records from the database, using upper
func (t *$MODELNAME$) GetAll(condition up.Cond) ([]*$MODELNAME$, error) {
    collection := upper.Collection(t.Table())
    var all []*$MODELNAME$

    res := collection.Find(condition)
    err := res.All(&all)
    if err != nil {
        return nil, err
    }

    return all, err
}

// Get gets one record from the database, by id, using upper
func (t *$MODELNAME$) Get$MODELNAME$ByID(id int) (*$MODELNAME$, error) {
    var one $MODELNAME$
    collection := upper.Collection(t.Table())

    res := collection.Find(up.Cond{"id": id})
    err := res.One(&one)
    if err != nil {
        return nil, err
    }
    return &one, nil
}

// Update updates a record in the database, using upper
func (t *$MODELNAME$) Update$MODELNAME$(m $MODELNAME$) error {
    m.UpdatedAt = time.Now().String()
    collection := upper.Collection(t.Table())
    res := collection.Find(m.ID)
    err := res.Update(&m)
    if err != nil {
        return err
    }
    return nil
}

// Delete deletes a record from the database by id, using upper
func (t *$MODELNAME$) Delete$MODELNAME$(id int) error {
    collection := upper.Collection(t.Table())
    res := collection.Find(id)
    err := res.Delete()
    if err != nil {
        return err
    }
    return nil
}

// Builder is an example of using upper's sql builder
func (t *$MODELNAME$) Builder(id int) ([]*$MODELNAME$, error) {
    collection := upper.Collection(t.Table())

    var result []*$MODELNAME$

    err := collection.Session().
        SQL().
        SelectFrom(t.Table()).
        Where("id > ?", id).
        OrderBy("id").
        All(&result)
    if err != nil {
        return nil, err
    }
    return result, nil
}
