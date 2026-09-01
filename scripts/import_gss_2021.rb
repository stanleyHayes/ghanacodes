#!/usr/bin/env ruby
require "digest"; require "json"; require "open3"; require "pathname"
EXPECTED_SHA256 = "06bb5ffc8cfeb59df612a10432209de82ddd8988d4e2fbb8c75e359dc2968560"
REGIONS = [["01","Western","WR"],["02","Central","CR"],["03","Greater Accra","GAR"],["04","Volta","VR"],["05","Eastern","ER"],["06","Ashanti","AR"],["07","Western North","WNR"],["08","Ahafo","AHR"],["09","Bono","BR"],["10","Bono East","BER"],["11","Oti","OR"],["12","Northern","NR"],["13","Savannah","SR"],["14","North East","NER"],["15","Upper East","UER"],["16","Upper West","UWR"]].freeze
NAME_OVERRIDES = {"0201"=>"Komenda Edina Eguafo Abirem Municipal"}.freeze
source_path=Pathname.new(ARGV.fetch(0){abort "usage: ruby scripts/import_gss_2021.rb <official-pdf> [output-json]"}).expand_path
output_path=Pathname.new(ARGV[1]||Pathname.new(__dir__).parent.join("data/codes.json").to_s).expand_path
abort "source checksum mismatch" unless Digest::SHA256.file(source_path).hexdigest==EXPECTED_SHA256
text,status=Open3.capture2("pdftotext","-layout",source_path.to_s,"-");abort "pdftotext failed" unless status.success?
start_at=text.rindex("APPENDIX 2: DISTRICT CODES");end_at=text.rindex("APPENDIX 3: COUNTRY CODES");abort "district appendix not found" unless start_at&&end_at&&end_at>start_at
region_by_abbreviation=REGIONS.to_h{|code,name,abbreviation|[abbreviation,{"code"=>code,"name"=>name}]}
districts=text[start_at...end_at].lines.grep(/\s[123]\s+\d{4}\s*$/).map do|line|
 abbreviation=line[/\A\S+/];suffix=line.match(/([123])\s+(\d{4})\s*$/);region=region_by_abbreviation.fetch(abbreviation)
 body=line.sub(/\A\S+\s+/,"").sub(/\s+[123]\s+\d{4}\s*$/,"");parsed_name=body.split(/\s{2,}/,2).first.to_s.strip.gsub(/\s+/," ");name=NAME_OVERRIDES.fetch(suffix[2],parsed_name)
 {"namespace"=>"gss-phc-2021-district","code"=>suffix[2],"name"=>name,"entityId"=>"geo:district:#{suffix[2]}","entityType"=>"district","districtType"=>{"1"=>"district","2"=>"municipal","3"=>"metropolitan"}.fetch(suffix[1]),"regionCode"=>region.fetch("code"),"validFrom"=>"2021-06-27","validTo"=>nil,"status"=>"active","sourceId"=>"gss-2021-phc-field-manual"}
end
abort "expected 261 district codes, found #{districts.length}" unless districts.length==261
abort "district codes are not unique" unless districts.map{|row|row.fetch("code")}.uniq.length==districts.length
abort "empty district name" if districts.any?{|row|row.fetch("name").empty?}
regions=REGIONS.map{|code,name,abbreviation|{"namespace"=>"gss-phc-2021-region","code"=>code,"name"=>name,"entityId"=>"geo:region:#{code}","entityType"=>"region","validFrom"=>"2021-06-27","validTo"=>nil,"status"=>"active","sourceId"=>"gss-2021-phc-field-manual","aliases"=>[abbreviation]}}
abbreviations=REGIONS.map{|code,name,abbreviation|{"namespace"=>"gss-phc-2021-region-abbreviation","code"=>abbreviation,"name"=>name,"entityId"=>"geo:region:#{code}","entityType"=>"region","validFrom"=>"2021-06-27","validTo"=>nil,"status"=>"active","sourceId"=>"gss-2021-phc-field-manual"}}
dataset={"version"=>"gss-phc-2021.1","effectiveAt"=>"2021-06-27","generatedFromSha256"=>EXPECTED_SHA256,"namespaces"=>[{"id"=>"gss-phc-2021-region","name"=>"GSS 2021 PHC region code","authority"=>"Ghana Statistical Service","entityType"=>"region"},{"id"=>"gss-phc-2021-region-abbreviation","name"=>"GSS 2021 PHC region abbreviation","authority"=>"Ghana Statistical Service","entityType"=>"region"},{"id"=>"gss-phc-2021-district","name"=>"GSS 2021 PHC district code","authority"=>"Ghana Statistical Service","entityType"=>"district"}],"codes"=>regions+abbreviations+districts}
output_path.dirname.mkpath;output_path.write(JSON.pretty_generate(dataset)+"\n");puts "Generated #{dataset.fetch('codes').length} codes across #{dataset.fetch('namespaces').length} namespaces"
