
import Std.Data.HashMap
import Std.Data.HashSet

open Std

abbrev Cid := String

structure Block where
  cid     : Cid
  accept  : List Cid
  reject  : List Cid

abbrev BlockMap := HashMap Cid Block
abbrev RejectionContext := HashMap Cid (HashSet Cid)

partial def computeVisibleView
  (blocks : BlockMap)
  (trustedTips : List Cid)
  : HashSet Cid :=
  let rec visit
    (stack : List Cid)
    (visited : HashSet Cid)
    (visible : HashSet Cid)
    (context : RejectionContext)
    : HashSet Cid :=
      match stack with
      | [] => visible
      | cid :: rest =>
        if visited.contains cid then
          visit rest visited visible context
        else 
          match blocks.get? cid with
          | none => visit rest (visited.insert cid) visible context
          | some block =>
            let inheritedRej :=
              blocks.fold (init := ∅) fun acc _ child =>
                if child.accept.contains cid then
                  acc ∪ context.getD child.cid ∅
                else acc
            let fullRej := inheritedRej ∪ HashSet.ofList block.reject
            let context' := context.insert cid fullRej
            if fullRej.contains cid then
              visit rest (visited.insert cid) visible context'
            else
              let visited' := visited.insert cid
              let visible' := visible.insert cid
              let stack' := block.accept ++ rest
              visit stack' visited' visible' context'
  visit trustedTips ∅ ∅ (HashMap.emptyWithCapacity 10)


def makeBlockMap : List Block → BlockMap :=
  HashMap.ofList ∘ (List.map fun b => Prod.mk b.cid b)


private def ex0 : BlockMap :=
  makeBlockMap
    [
      ⟨ "B0" , ∅ , ∅ ⟩ ,
    ]
#eval (computeVisibleView ex0 ["B0"]).toList = ["B0"]

private def ex1 : BlockMap :=
  makeBlockMap
    [
      ⟨ "B0" , ∅ , ∅ ⟩ ,
      ⟨ "B1" , ["B0"] , ∅ ⟩ ,
    ]
#eval (computeVisibleView ex1 ["B1"]).toList = ["B0" , "B1"]

private def ex2 : BlockMap :=
  makeBlockMap
    [
      ⟨ "B0" , ∅ , ∅ ⟩ ,
      ⟨ "B1" , ["B0"] , ∅ ⟩ ,
      ⟨ "B2" , ["B1"] , ∅ ⟩ ,
    ]
#eval (computeVisibleView ex2 ["B2"]).toList = ["B0", "B2", "B1"]

private def ex3 : BlockMap :=
  makeBlockMap
    [
      ⟨ "B0" , ∅ , ∅ ⟩ ,
      ⟨ "B1" , ["B0"] , ∅ ⟩ ,
      ⟨ "B2" , ["B1"] , ∅ ⟩ ,
      ⟨ "B3" , ["B2"] , ["B0"] ⟩ ,
    ]
#eval (computeVisibleView ex3 ["B3"]).toList = ["B2", "B3", "B1"]

private def ex4 : BlockMap :=
  makeBlockMap
    [
      ⟨ "B0" , ∅ , ∅ ⟩ ,
      ⟨ "B1" , ["B0"] , ∅ ⟩ ,
      ⟨ "B2" , ["B1"] , ∅ ⟩ ,
      ⟨ "B4" , ["B2"] , ["B1"] ⟩ ,
    ]
#eval (computeVisibleView ex4 ["B4"]).toList = ["B4", "B2"]

private def ex5 : BlockMap :=
  makeBlockMap
    [
      ⟨ "B0" , ∅ , ∅ ⟩ ,
      ⟨ "B1" , ["B0"] , ∅ ⟩ ,
      ⟨ "B2" , ["B1"] , ∅ ⟩ ,
      ⟨ "B5" , ["B2"] , ["B2"] ⟩ ,
    ]
#eval (computeVisibleView ex5 ["B5"]).toList = ["B5"]

private def ex6 : BlockMap :=
  makeBlockMap
    [
      ⟨ "B0" , ∅ , ∅ ⟩ ,
      ⟨ "B1" , ["B0"] , ∅ ⟩ ,
      ⟨ "B2" , ["B1"] , ∅ ⟩ ,
      ⟨ "B6" , ["B1"] , ∅ ⟩ ,
      ⟨ "B7" , ["B2", "B6"] , ∅ ⟩ ,
    ]
#eval (computeVisibleView ex6 ["B7"]).toList = ["B0", "B2", "B7", "B1", "B6"]

private def ex7 : BlockMap :=
  makeBlockMap
    [
      ⟨ "B0" , ∅ , ∅ ⟩ ,
      ⟨ "B1" , ["B0"] , ∅ ⟩ ,
      ⟨ "B2" , ["B1"] , ∅ ⟩ ,
      ⟨ "B6" , ["B1"] , ∅ ⟩ ,
      ⟨ "B7" , ["B2", "B6"] , ∅ ⟩ ,
      ⟨ "B8" , ["B7"] , ["B6"] ⟩ ,
    ]
#eval (computeVisibleView ex7 ["B8"]).toList = ["B0", "B2", "B7", "B1", "B8"]

private def ex8 : BlockMap :=
  makeBlockMap
    [
      ⟨ "B0" , ∅ , ∅ ⟩ ,
      ⟨ "B1" , ["B0"] , ∅ ⟩ ,
      ⟨ "B2" , ["B1"] , ∅ ⟩ ,
      ⟨ "B6" , ["B1"] , ∅ ⟩ ,
      ⟨ "B7" , ["B2", "B6"] , ∅ ⟩ ,
      ⟨ "B8" , ["B7"] , ["B6"] ⟩ ,
      ⟨ "B9" , ["B1"] , ∅ ⟩ ,
      ⟨ "B10" , ["B9"] , ∅ ⟩ ,
      ⟨ "B11" , ["B8", "B10"] , ∅ ⟩ ,
    ]
#eval (computeVisibleView ex8 ["B11"]).toList = ["B0", "B10", "B2", "B7", "B11", "B1", "B8", "B9"]

private def ex9 : BlockMap :=
  makeBlockMap
    [
      ⟨ "B0" , ∅ , ∅ ⟩ ,
      ⟨ "B1" , ["B0"] , ∅ ⟩ ,
      ⟨ "B2" , ["B1"] , ∅ ⟩ ,
      ⟨ "B6" , ["B1"] , ∅ ⟩ ,
      ⟨ "B7" , ["B2", "B6"] , ∅ ⟩ ,
      ⟨ "B8" , ["B7"] , ["B6"] ⟩ ,
      ⟨ "B9" , ["B1"] , ∅ ⟩ ,
      ⟨ "B10" , ["B9"] , ∅ ⟩ ,
      ⟨ "B11" , ["B8", "B10"] , ∅ ⟩ ,
      ⟨ "B12" , ["B11"] , ["B2"] ⟩ ,
    ]
#eval (computeVisibleView ex9 ["B12"]).toList = ["B0", "B10", "B7", "B11", "B12", "B1", "B8", "B9"]

private def ex10 : BlockMap :=
  makeBlockMap
    [
      ⟨ "B0" , ∅ , ∅ ⟩ ,
      ⟨ "B1" , ["B0"] , ∅ ⟩ ,
      ⟨ "B2" , ["B1"] , ∅ ⟩ ,
      ⟨ "B6" , ["B1"] , ∅ ⟩ ,
      ⟨ "B7" , ["B2", "B6"] , ∅ ⟩ ,
      ⟨ "B8" , ["B7"] , ["B6"] ⟩ ,
      ⟨ "B9" , ["B1"] , ∅ ⟩ ,
      ⟨ "B10" , ["B9"] , ∅ ⟩ ,
      ⟨ "B11" , ["B8", "B10"] , ∅ ⟩ ,
      ⟨ "B12" , ["B11"] , ["B2"] ⟩ ,
      ⟨ "B13" , ["B7", "B12"] , ∅ ⟩ ,
    ]
#eval (computeVisibleView ex10 ["B13"]).toList = ["B0", "B10", "B11", "B12", "B2", "B7", "B13", "B1", "B8", "B6", "B9"]

private def ex11 : BlockMap :=
  makeBlockMap
    [
      ⟨ "B0" , ∅ , ∅ ⟩ ,
      ⟨ "B1" , ["B0"] , ∅ ⟩ ,
      ⟨ "B2" , ["B1"] , ∅ ⟩ ,
      ⟨ "B6" , ["B1"] , ∅ ⟩ ,
      ⟨ "B7" , ["B2", "B6"] , ∅ ⟩ ,
      ⟨ "B8" , ["B7"] , ["B6"] ⟩ ,
      ⟨ "B9" , ["B1"] , ∅ ⟩ ,
      ⟨ "B10" , ["B9"] , ∅ ⟩ ,
      ⟨ "B11" , ["B8", "B10"] , ∅ ⟩ ,
      ⟨ "B12" , ["B11"] , ["B2"] ⟩ ,
      ⟨ "B13" , ["B7", "B12"] , ∅ ⟩ ,
      ⟨ "B14" , ["B13"] , ["B9"] ⟩ ,
    ]
#eval (computeVisibleView ex11 ["B14"]).toList = ["B0", "B10", "B11", "B12", "B2", "B7", "B13", "B1", "B8", "B6", "B14"]

private def ex12 : BlockMap :=
  makeBlockMap
    [
      ⟨ "A" , ∅ , ∅ ⟩ ,
      ⟨ "B" , ["A"] , ∅ ⟩ ,
      ⟨ "C" , ["A"] , ∅ ⟩ ,
      ⟨ "D" , ["C"] , ∅ ⟩ ,
      ⟨ "E" , ["D"] , ∅ ⟩ ,
      ⟨ "F" , ["B", "E"] , ["D"] ⟩ ,
      ⟨ "G" , ["F"] , ∅ ⟩ ,
    ]
#eval (computeVisibleView ex12 ["G"]).toList = ["B", "G", "A", "E", "F"]

private def ex13 : BlockMap :=
  makeBlockMap
    [
      ⟨ "A" , ∅ , ∅ ⟩ ,
      ⟨ "B" , ["A"] , ["C"] ⟩ ,
      ⟨ "C" , ["A"] , ∅ ⟩ ,
      ⟨ "D" , ["C"] , ∅ ⟩ ,
      ⟨ "E" , ["D"] , ∅ ⟩ ,
      ⟨ "F" , ["B", "E"] , ∅ ⟩ ,
      ⟨ "G" , ["F"] , ∅ ⟩ ,
    ]
#eval (computeVisibleView ex13 ["G"]).toList = ["B", "G", "D", "A", "C", "E", "F"]

def ex14 : BlockMap :=
  makeBlockMap
    [
      ⟨ "R" , ∅ , ∅ ⟩,
      ⟨ "B" , ["R"] , ["R"] ⟩,
      ⟨ "A" , ["B"] , ∅ ⟩,
      ⟨ "Y" , ["B"] , ∅ ⟩,
      ⟨ "X" , ["Y"] , ∅ ⟩,
      ⟨ "T" , ["A", "X"] , ∅ ⟩,
    ]
#eval (computeVisibleView ex14 ["T"]).toList = ["X", "B", "T", "A", "Y"]
