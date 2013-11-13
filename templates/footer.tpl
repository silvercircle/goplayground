{{define "footer"}}
</div> <!-- main_content_section -->
</div> <!-- content_section -->
<div class="clear" id="footer_section">
  <div>
  </div>
  <div class="righttext floatright">{{.loadtime}}</div>
  <div class="copyright">
    <span>Route: {{.matched_route}}</span>
  </div>
  <div>
    <a id="button_xhtml" href="http://validator.w3.org/check?uri=referer" target="_blank" class="new_win" title="Valid HTML"><span>HTML</span></a> |
  </div>
</div>
</div> <!-- #wrap -->
</body>
</html>
{{end}}

