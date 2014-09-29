package main

type Remote interface {
	Put(PutRequest) PutResponse
	Get(GetRequest) GetResponse
	Delete(DeleteRequest) DeleteResponse
	Head(HeadRequest) HeadResponse
	List(ListRequest) ListResponse
}
