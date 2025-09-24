#include<bits/stdc++.h>
using namespace std;

class Animal{
    public:
    string s;
    Animal(string str):s(str){}
    virtual void call(){
        cout<<"k"<<endl;
    }
    virtual void Call() const = 0;
};

class Cat: public Animal{
    public:
    Cat(string str):Animal(str){}
    void call(){
        cout<<"miao"<<endl;
    }
    void Call() const override{
        cout<<"miao"<<endl;
    }
};

class Dog: public Animal{
    public:
    Dog(string str):Animal(str){}
    void call(){
        cout<<"wang"<<endl;
    }
    void Call() const override{
        cout<<"wang"<<endl;
    }
};

int main(){
    Animal *a =new  Cat("a");
    a->call();
    a->Call();
    a = new Dog("dog");
    a->call();
    a->Call();
    return 0;
}