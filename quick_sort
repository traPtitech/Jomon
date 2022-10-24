#include <iostream>
#include <algorithm>
using namespace std;
 
// Partition using the Lomuto partition scheme
int partition(int a[], int start, int end)
{
    // Pick the rightmost element as a pivot from the array
    int pivot = a[end];
 
    // elements less than the pivot will be pushed to the left of `pIndex`
    // elements more than the pivot will be pushed to the right of `pIndex`
    // equal elements can go either way
    int pIndex = start;
 
    // each time we find an element less than or equal to the pivot, `pIndex`
    // is incremented, and that element would be placed before the pivot.
    for (int i = start; i < end; i++)
    {
        if (a[i] <= pivot)
        {
            swap(a[i], a[pIndex]);
            pIndex++;
        }
    }
 
    // swap `pIndex` with pivot
    swap (a[pIndex], a[end]);
 
    // return `pIndex` (index of the pivot element)
    return pIndex;
}
 
// Quicksort routine
void quicksort(int a[], int start, int end)
{
    // base condition
    if (start >= end) {
        return;
    }
 
    // rearrange elements across pivot
    int pivot = partition(a, start, end);
 
    // recur on subarray containing elements that are less than the pivot
    quicksort(a, start, pivot - 1);
 
    // recur on subarray containing elements that are more than the pivot
    quicksort(a, pivot + 1, end);
}
 
// C++ implementation of the Quicksort algorithm
int main()
{
    int a[] = { 9, -3, 5, 2, 6, 8, -6, 1, 3 };
    int n = sizeof(a)/sizeof(a[0]);
 
    quicksort(a, 0, n - 1);
 
    // print the sorted array
    for (int i = 0; i < n; i++) {
        cout << a[i] << " ";
    }
 
    return 0;
}
