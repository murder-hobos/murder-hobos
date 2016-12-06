function Sort() {

    var input, filter, set, li, a, i;
    input = document.getElementById("mySearch");
    filter = input.value.toUpperCase();
    set = document.getElementById("Spell_List");
    li = set.getElementsByTagName("Spell");


    

    // Loops through query and narrows results, not working rn
    
        for (i = 0; i < li.length; i++) {
            a = li[i].getElementsByTagName("")[0];
            if (a.innerHTML.toUpperCase().indexOf(filter) > -1) {
                li[i].style.display = "";
            } else {
                li[i].style.display = "none";
            }
        }


       
}
/*var searchGen = {
  valueNames: [ 'Spells']
};

var userList = new List('SpellsSearched', searchGen);*/

/*var $ul = $("ul.Spell_List");
var searchParameters = document.getElementById("mySearch").text;*/

