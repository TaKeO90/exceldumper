# exceltovcf

 * *WARNING*  Program still under devlopement so it may be some bugs that should i fix later ... don't hesitate to submit an issue .

 * The program is Just a simple tool to transform excel data into a csv or vcf file cause i need that some times when i have a huge excel file with a lot of contact 
    numbers and i want a single csv or vcf file to import it into my phone and have those concat numbers on it .

# Installation 

 ```sh
    $ git clone https://github.com/TaKeO90/exceldumper
    $ make -B `cli program`
    $ make server `back-end web server api`
 ```

# USAGE 


```sh
    * excel to vcf
    $ ./exceldumper -excel <excel file name> -sheet <sheet name>  -cntnumber <number of rows to dump> -vcffile <output vcf file name>
    * excel to csv
    $ ./exceldumper -excel <excel file name> -sheet <sheet name> -csvfile <output csv file name>
```

