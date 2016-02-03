(function($){
  function resize(windowHeight){
    $('.promo').css('min-height', windowHeight - $('header').height());
    $('.showoff, .features').css('min-height', windowHeight);
  }

  $(document).ready(function() {
    var windowHeight = $(window).height();

    resize(windowHeight);

    $('.run').click(function() {
      $('.embed').addClass('show');
    });

  });

  $(window).resize(function() {
    var windowHeight = $(window).height();
    
    resize(windowHeight);
  });

  $(window).scroll(function() {
    var topPosition = $(window).scrollTop();

    if (topPosition >= $('.promo').height()) {
      // $('body').css('background','#11C9C3');
      $('.homePage').css('background','#11C9C3');
      
      if (topPosition >= $('.promo').height()*2) {
        $('.homePage').css('background','#FF836E');
      }
    } else {
      $('.homePage').removeAttr('style');
    }
  });

  

})(jQuery);