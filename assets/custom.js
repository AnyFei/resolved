
/* jshint ignore:end */

//$('.collapse').collapse();
function makeEditable(){
  $("#saveChangesBtn").toggle();
  $(".editable").each(function() {  
      if($(this).attr("contenteditable") != "true"){
        $(this).attr("contenteditable", "true");
        $("#editItemBtn").html("Cancel");
        $(this).toggleClass( "editElementActive" );
        $(this).on('DOMSubtreeModified', $(this), function() {
          key = $(this).attr('id');
          value = $(this).html();
          updatedValues[key] = value;
          console.log("testing");
      });
      }
     else{
        $("#editItemBtn").html("Edit");
        $(this).attr("contenteditable", "false");
        $(this).toggleClass( "editElementActive" );


    }
    });
}
function deleteItem(itemId, itemType){
  const requestOptions = {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ID: itemId, type: itemType}),
  };
  fetch("/delete", requestOptions)
  .then(function(){ // expecting a json response
  window.location.href = "/"+itemType.toLowerCase()+"s/";
}  );
        }

        document.addEventListener("DOMContentLoaded", function(event) { 
          const requestOptions = {
          method: "GET",
          headers: { "Content-Type": "application/json" },
          };
          fetch("/refresh_token", requestOptions)
          .then(function(){ // expecting a json response
            console.log("token refreshed")
          });
        
      });

function saveChanges(type, id){
  updatedValues["type"] = type;
  updatedValues["id"] = id;
    const requestOptions = {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(updatedValues),
  };
  fetch("/update", requestOptions)
  .then(makeEditable());}


    function showRemoveModal(){
    // Get the modal
var modal = document.getElementById("myModal");

// Get the button that opens the modal
var btn = document.getElementById("myBtn");

// Get the <span> element that closes the modal
var cancelBtn = document.getElementById("cancelRemove");

// When the user clicks on the button, open the modal
btn.onclick = function() {
  modal.style.display = "block";
}

// When the user clicks on <span> (x), close the modal
  cancelBtn.onclick = function() {
  modal.style.display = "none";
}

// When the user clicks anywhere outside of the modal, close it
window.onclick = function(event) {
  if (event.target == modal) {
    modal.style.display = "none";
  }
}
}

function updatePermissions(permissionType, userId){
  let currentStatus = $("#"+permissionType).hasClass("permission-true");
  currentStatus = !currentStatus;
  console.log(currentStatus);
  const requestOptions = {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({type: "UserPermissions", permission: permissionType, newValue: currentStatus.toString(), userId: userId}),
  };
  fetch("/update", requestOptions)
  .then(function(){
    if(currentStatus){
      $("#"+permissionType).removeClass("permission-false");
      $("#"+permissionType).addClass("permission-true");
    }
    else{
      $("#"+permissionType).removeClass("permission-true");
    }
  }
  );

}


